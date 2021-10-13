package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Status      string             `bson:"status" json:"status"`
	Message     string             `bson:"message" json:"message"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"updatedDate"`
	Total       float64            `bson:"total" json:"total"`
	TotalCost   float64            `bson:"totalCost" json:"totalCost"`
	Type        string             `bson:"type" json:"type"`
}

type OrderDetail struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Status      string             `bson:"status" json:"status"`
	Message     string             `bson:"message" json:"message"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"updatedDate"`
	Total       float64            `bson:"total" json:"total"`
	TotalCost   float64            `bson:"totalCost" json:"totalCost"`
	Type        string             `bson:"type" json:"type"`
	Items       []OrderItemDetail  `json:"items"`
	Payment     Payment            `json:"payment"`
}
