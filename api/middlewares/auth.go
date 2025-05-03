package middlewares

import (
	"context"
	"net/http"

	"github.com/arcedo/financial-ai-backend/api/helpers"
	db "github.com/arcedo/financial-ai-backend/database"
	"github.com/arcedo/financial-ai-backend/types"
	"github.com/arcedo/financial-ai-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func JWTAuthMiddleware(secretKey []byte, store db.MongoStorage) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				errValue := utils.ErrorMap[utils.ErrUnauthorized]
				helpers.WriteJSON(w, http.StatusUnauthorized, nil, &errValue, "")
				return
			}

			claims, err := utils.ValidateJWT(secretKey, tokenString)
			if err != nil {
				helpers.WriteJSON(w, http.StatusUnauthorized, nil, &utils.APIError{
					Code:    "UNAUTHORIZED",
					Message: err.Error(),
				}, "")
				return
			}

			// Convert string ID back to ObjectID
			userID, err := primitive.ObjectIDFromHex(claims.ID)
			if err != nil {
				helpers.WriteJSON(w, http.StatusUnauthorized, nil, &utils.APIError{
					Code:    "INVALID_OBJECTID",
					Message: "invalid user ID in token",
				}, "")
				return
			}

			// Check if the user exists in the database
			userCollection := store.Collection("users")
			var foundUser types.User
			err = userCollection.FindOne(context.Background(), bson.D{
				{Key: "_id", Value: userID},
			}).Decode(&foundUser)

			if err != nil {
				if err == mongo.ErrNoDocuments {
					helpers.WriteJSON(w, http.StatusUnauthorized, nil, &utils.APIError{
						Code:    "USER_NOT_FOUND",
						Message: "user not found",
					}, "")
					return
				}
				helpers.WriteJSON(w, http.StatusInternalServerError, nil, &utils.APIError{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "error searching for user",
				}, "")
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
