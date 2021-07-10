package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Subscription struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	UserId      primitive.ObjectID `bson:"userId" json:"userId"`
	Channel     string             `bson:"channel" json:"channel"`
	DeviceToken string             `bson:"deviceToken" json:"deviceToken"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"updatedDate"`
}
