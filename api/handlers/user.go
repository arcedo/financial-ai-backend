package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/arcedo/financial-ai-backend/api/helpers"
	db "github.com/arcedo/financial-ai-backend/database"
	"github.com/arcedo/financial-ai-backend/requests"
	"github.com/arcedo/financial-ai-backend/types"
	"github.com/arcedo/financial-ai-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Login(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	var user types.User
	// Decode the incoming request body to extract the user credentials (e.g., username/password)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	if err := user.ValidateUserLogin(); err != nil {
		return err
	}

	userCollection := store.Collection("users")
	var foundUser types.User

	err := userCollection.FindOne(context.Background(), bson.D{
		{Key: "email", Value: user.Email},
	}).Decode(&foundUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("invalid credentials")
		}
		return fmt.Errorf("error searching for user: %v", err)
	}
	if passOk := utils.CheckPasswordHash(foundUser.Password, user.Password); passOk == false {
		return fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT([]byte(os.Getenv("SECRET")), foundUser.ID)
	if err != nil {
		return err
	}

	publicUser := types.PublicUser{
		Name:     foundUser.Name,
		LastName: foundUser.LastName,
		Email:    foundUser.Email,
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  publicUser,
	}, nil, "")
	return nil
}

func CreateUser(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	var newUser = types.NewUser{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		return err
	}

	if err := newUser.ValidateUserCreation(); err != nil {
		return err
	}

	// Check if user with the same email already exists
	userCollection := store.Collection("users")
	var existingUser types.PublicUser
	err := userCollection.FindOne(context.Background(), bson.D{
		{Key: "email", Value: newUser.Email},
	}).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking for existing user: %v", err)
	}
	if existingUser.Email != "" {
		return fmt.Errorf("user with this email already exists")
	}

	// Hash the password before storing it
	newUser.Password, err = utils.HashPassword(newUser.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	newUser.Name = utils.SanitizeString(newUser.Name)
	newUser.LastName = utils.SanitizeString(newUser.LastName)
	newUser.Email = utils.SanitizeString(newUser.Email)

	// Insert the new user into the database
	res, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("error inserting new user: %v", err)
	}

	// Retrieve the inserted ID and verify it's an ObjectID
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unexpected ID type")
	}

	// Generate the JWT token using the user ID
	token, err := utils.GenerateJWT([]byte(os.Getenv("SECRET")), insertedID)
	if err != nil {
		return fmt.Errorf("error generating token: %v", err)
	}

	publicUser := types.PublicUser{
		Name:     newUser.Name,
		LastName: newUser.LastName,
		Email:    newUser.Email,
	}

	// Return the response with the generated token
	helpers.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"token": token,
		"user":  publicUser,
	}, nil, "user created successfully")
	return nil
}

func GetUser(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unable to retrieve user ID from context")
	}

	// Now you can use userID directly in the query
	userCollection := store.Collection("users")
	var user types.PublicUser
	err := userCollection.FindOne(context.Background(), bson.D{
		{Key: "_id", Value: userID},
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("error retrieving user: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, user, nil, "")
	return nil
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	userCollection := store.Collection("users")
	cursor, err := userCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return fmt.Errorf("error retrieving users: %v", err)
	}
	defer cursor.Close(context.Background())

	var users []types.User
	for cursor.Next(context.Background()) {
		var user types.User
		if err := cursor.Decode(&user); err != nil {
			return fmt.Errorf("error decoding user: %v", err)
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor error: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, users, nil, "")
	return nil
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unable to retrieve user ID from context")
	}

	userCollection := store.Collection("users")
	var userData types.Insights
	err := userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&userData.User)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("error retrieving user: %v", err)
	}

	transactionsColl := store.Collection("transactions")
	cursor, err := transactionsColl.Find(r.Context(), bson.M{"user_id": userID})
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer cursor.Close(r.Context())

	if err := cursor.All(r.Context(), &userData.Transactions); err != nil {
		return fmt.Errorf("failed to decode transactions: %v", err)
	}

	// Calculate position summary
	var position types.PositionSummary
	for _, t := range userData.Transactions {
		switch t.Type {
		case "buy":
			position.TotalBuys += t.Amount
		case "sell":
			position.TotalSells += t.Amount
		case "save":
			position.TotalSaves += t.Amount
		case "entry":
			position.TotalEntry += t.Amount
		}
	}
	position.NetMarket = position.TotalBuys - position.TotalSells
	position.NetBalance = position.NetMarket + position.TotalSaves + position.TotalEntry
	userData.Position = position

	// Send to LLM
	newProfile, err := requests.RequestUpdateUserProfile(userData)
	if err != nil {
		return fmt.Errorf("failed to update LLM profile: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, newProfile, nil, "User profile updated successfully")
	return nil
}

func GetRecommendations(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unable to retrieve user ID from context")
	}

	recommendations, err := requests.GetRecommendations(userID.Hex())
	if err != nil {
		return fmt.Errorf("failed to get recommendations: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, recommendations, nil, "")
	return nil
}

func GetAdvice(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unable to retrieve user ID from context")
	}

	advice, err := requests.GetAdvice(userID.Hex())
	if err != nil {
		return fmt.Errorf("failed to get advice: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, advice, nil, "")
	return nil
}

func GetAssetRecommendation(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unable to retrieve user ID from context")
	}

	symbol, err := utils.GetPathParam(r.URL.Path, 1)
	if err != nil {
		return fmt.Errorf("unable to retrieve symbol from URL: %v", err)
	}

	recommendations, err := requests.GetAssetRecommendation(userID.Hex(), symbol)
	if err != nil {
		return fmt.Errorf("failed to get recommendations: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, recommendations, nil, "")
	return nil
}
