package product

import (
	"devper/app/featues/product/repository"
	"devper/app/featues/product/usecase"
	"devper/middlewares"
	"devper/utils/constant"
	"github.com/gin-gonic/gin"
)

func ApplyProductAPI(app *gin.RouterGroup, productEntity repository.IProduct) {

	productRoute := app.Group("product")

	productRoute.GET("",
		usecase.GetProducts(productEntity),
	)

	productRoute.POST("",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.CreateProduct(productEntity),
	)

	productRoute.GET("/:productId",
		usecase.GetProductById(productEntity),
	)

	productRoute.PUT("/:productId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.UpdateProductById(productEntity),
	)

	productRoute.DELETE("/:productId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.DeleteProductById(productEntity),
	)

	productRoute.GET("/serial-number/:serialNumber",
		usecase.GetProductBySerialNumber(productEntity),
	)

	productRoute.GET("/:productId/lot",
		usecase.GetLotsByProductId(productEntity),
	)

	productRoute.PUT("/lot/:lotId",
		usecase.UpdateLotById(productEntity),
	)

	productRoute.GET("/lot/:lotId",
		usecase.GetLotById(productEntity),
	)
}
