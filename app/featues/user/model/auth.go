package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserReference struct {
	Id          primitive.ObjectID `bson:"_id" json:"userRefId"`
	UserId      primitive.ObjectID `bson:"userId" json:"-"`
	Type        string             `bson:"type" json:"type"`
	Objective   string             `bson:"objective" json:"objective"`
	Channel     string             `bson:"channel" json:"channel"`
	ChannelInfo string             `bson:"channelInfo" json:"channelInfo"`
	RefId       string             `bson:"refId" json:"refId"`
	Code        string             `bson:"code" json:"-"`
	Status      string             `bson:"status" json:"status"`
	ValidPeriod int                `bson:"validPeriod" json:"validPeriod"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	ExpireDate  time.Time          `bson:"expireDate" json:"expireDate"`
}
