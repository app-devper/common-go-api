package form

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUser struct {
	Username  string `json:"username" binding:"required"`
	Objective string `json:"objective" binding:"required"`
}

type Channel struct {
	Channel     string `json:"channel"`
	ChannelInfo string `json:"channelInfo"`
}

type VerifyChannel struct {
	UserRefId   string `json:"userRefId" binding:"required"`
	Channel     string `json:"channel" binding:"required"`
	ChannelInfo string `json:"channelInfo" binding:"required"`
}

type VerifyCode struct {
	UserRefId string `json:"userRefId" binding:"required"`
	RefId     string `json:"refId" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

type Reference struct {
	UserId      primitive.ObjectID
	Type        string
	Objective   string
	Channel     string
	ChannelInfo string
	Status      string
	ValidPeriod int
	ExpireDate  time.Time
}
