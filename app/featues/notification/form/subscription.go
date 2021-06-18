package form

type Subscription struct {
	UserId      string `json:"userId" binding:"required"`
	Channel     string `json:"channel" binding:"required"`
	DeviceToken string `json:"deviceToken" binding:"required"`
}
