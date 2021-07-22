package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderItem struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	OrderId     primitive.ObjectID `bson:"orderId" json:"orderId"`
	ProductId   primitive.ObjectID `bson:"productId" json:"productId"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Price       float32            `bson:"price" json:"price"`
	Discount    float32            `bson:"discount" json:"discount"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"updatedDate"`
}
