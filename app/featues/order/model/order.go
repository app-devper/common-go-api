package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Status      string             `bson:"status" json:"status"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"updatedDate"`
}
