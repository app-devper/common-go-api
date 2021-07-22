package form

type OrderItem struct {
	ProductId string  `json:"productId" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required"`
	Price     float32 `json:"price" binding:"required"`
	Discount  float32 `json:"discount"`
}
