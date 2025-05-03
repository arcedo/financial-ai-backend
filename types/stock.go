package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Stock struct {
	ID         primitive.ObjectID `json:"_id"`
	Symbol     string             `json:"symbol"`
	Date       string             `json:"date"`
	OpenPrice  float32            `json:"open"`
	ClosePrice float32            `json:"close"`
	HighPrice  float32            `json:"high"`
	LowPrice   float32            `json:"low"`
	Volume     float32            `json:"volume"`
}

type NewStock struct {
	Symbol     string  `json:"symbol"`
	Date       string  `json:"date"`
	OpenPrice  float32 `json:"open"`
	ClosePrice float32 `json:"close"`
	HighPrice  float32 `json:"high"`
	LowPrice   float32 `json:"low"`
	Volume     float32 `json:"volume"`
}
