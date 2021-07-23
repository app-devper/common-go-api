package form

type Product struct {
	Name         string  `json:"name" binding:"required"`
	NameEn       string  `json:"nameEn"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" binding:"required"`
	Unit         string  `json:"unit"`
	Quantity     int     `json:"quantity" binding:"required"`
	SerialNumber string  `json:"serialNumber" binding:"required"`
	LotNumber    string  `json:"lotNumber"`
	ExpireDate   string  `json:"expireDate"`
}

type UpdateProduct struct {
	Name        string  `json:"name" binding:"required"`
	NameEn      string  `json:"nameEn"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Unit        string  `json:"unit"`
	Quantity    int     `json:"quantity" binding:"required"`
}

type ProductLot struct {
	Quantity   int    `json:"quantity" binding:"required"`
	LotNumber  string `json:"lotNumber" binding:"required"`
	ExpireDate string `json:"expireDate" binding:"required"`
}
