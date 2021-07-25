package form

type Notify struct {
	Message string `json:"message" binding:"required"`
}
