package types

import (
	"github.com/arcedo/financial-ai-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name"`
	LastName       string             `json:"last_name"`
	Email          string             `json:"email"`
	Password       string             `json:"password,omitempty"`
	RiskScore      int                `json:"risk_score"`
	FinancialScore int                `json:"financial_score"`
}

type UserProfile struct {
	RiskScore      int `json:"risk_score"`
	FinancialScore int `json:"financial_score"`
}

type PublicUser struct {
	Name           string `json:"name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	RiskScore      int    `json:"risk_score"`
	FinancialScore int    `json:"financial_score"`
}

type NewUser struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *NewUser) ValidateUserCreation() error {
	if err := utils.ValidateStringField(u.Name, "Name"); err != nil {
		return err
	}

	if err := utils.ValidateStringField(u.Name, "Last name"); err != nil {
		return err
	}

	if err := utils.ValidateEmail(u.Email); err != nil {
		return err
	}

	if err := utils.ValidatePassword(u.Password); err != nil {
		return err
	}

	return nil
}

func (u *User) ValidateUserLogin() error {
	if err := utils.ValidateStringField(u.Email, "email"); err != nil {
		return err
	}

	if err := utils.ValidateStringField(u.Password, "password"); err != nil {
		return err
	}

	return nil
}
