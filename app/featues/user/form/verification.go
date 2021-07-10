package form

type VerifyUser struct {
	Username  string `json:"username" binding:"required"`
	Objective string `json:"objective" binding:"required"`
}

type Channel struct {
	Channel     string `json:"channel"`
	ChannelInfo string `json:"channelInfo"`
}

type VerifyRequest struct {
	UserRefId   string `json:"userRefId" binding:"required"`
	Channel     string `json:"channel" binding:"required"`
	ChannelInfo string `json:"channelInfo" binding:"required"`
}

type VerifyCode struct {
	UserRefId string `json:"userRefId" binding:"required"`
	RefId     string `json:"refId" binding:"required"`
	Code      string `json:"code" binding:"required"`
}
