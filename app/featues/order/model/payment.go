package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Payment struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	OrderId     primitive.ObjectID `bson:"orderId" json:"orderId"`
	Status      string             `bson:"status" json:"status"`
	Amount      float32            `bson:"amount" json:"amount"`
	Type        string             `bson:"type" json:"type"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"updatedDate"`
}
