package form

type OrderItem struct {
	ProductId string  `json:"productId" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	CostPrice float64 `json:"costPrice"`
	Discount  float64 `json:"discount"`
}
