package types

import (
	"fmt"

	"github.com/arcedo/financial-ai-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID     int                `json:"_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	Type   string             `json:"type"`
	Amount float64            `json:"amount`
	Date   string             `json:"date"`
	Symbol string             `json:"symbol", bson:"symbol"`
}

type TransactionPublic struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
	Symbol string  `json:"symbol", bson:"symbol"`
}

type NewTransaction struct {
	Type   string             `json:"type"`
	Amount float64            `json:"amount"`
	Date   string             `json:"date"`
	Symbol string             `json:"symbol", bson:"symbol"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
}

func ValidateTransaction(transaction TransactionPublic) error {
	if err := utils.ValidateStringField(transaction.Type, "type"); err != nil {
		return err
	}

	if transaction.Type != "buy" && transaction.Type != "sell" && transaction.Type != "entry" && transaction.Type != "save" {
		return fmt.Errorf("invalid transaction type: %s, must be buy, sell, entry or save", transaction.Type)
	}

	if err := utils.ValidateStringField(transaction.Date, "date"); err != nil {
		return err
	}

	if err := utils.ValidateStringField(transaction.Symbol, "symbol"); err != nil && transaction.Type != "entry" && transaction.Type != "save" {
		return err
	}
	return nil
}
