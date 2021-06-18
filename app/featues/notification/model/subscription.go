package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscription struct {
	UserId      primitive.ObjectID `bson:"_id" json:"id"`
	Channel     string             `bson:"channel" json:"channel"`
	DeviceToken string             `bson:"deviceToken" json:"deviceToken"`
	CreatedDate string             `bson:"createdDate" json:"createdDate"`
}
