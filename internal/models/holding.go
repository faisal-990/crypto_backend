package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Holding struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID string             `bson:"user_id" json:"userId"`
	CoinID string             `bson:"coin_id" json:"coinId"`
	Amount float64            `bson:"amount" json:"amount"`
}

type Snapshot struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     string             `bson:"user_id" json:"userId"`
	TotalValue float64            `bson:"total_value" json:"totalValue"`
	Timestamp  primitive.DateTime `bson:"timestamp" json:"timestamp"`
}

func ToPrimitiveDateTime(t time.Time) primitive.DateTime {
	return primitive.NewDateTimeFromTime(t)
}
