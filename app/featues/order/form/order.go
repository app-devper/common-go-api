package form

type Order struct {
	Items  []OrderItem `json:"items" binding:"required"`
	Amount float32     `json:"amount" binding:"required"`
	Type   string      `json:"type" binding:"required"`
}
