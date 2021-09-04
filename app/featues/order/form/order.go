package form

import "time"

type Order struct {
	Items     []OrderItem `json:"items" binding:"required"`
	Amount    float64     `json:"amount" binding:"required"`
	Type      string      `json:"type" binding:"required"`
	Total     float64     `json:"total"`
	TotalCost float64     `json:"totalCost"`
	Change    float64     `json:"change"`
	Message   string      `json:"message"`
}

type GetOrderRange struct {
	StartDate time.Time `form:"startDate" binding:"required"`
	EndDate   time.Time `form:"endDate" binding:"required"`
}
