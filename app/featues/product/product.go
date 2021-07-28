package product

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/product/form"
	"mgo-gin/app/featues/product/repository"
	"mgo-gin/db"
	"mgo-gin/middlewares"
	"net/http"
)

func ApplyProductAPI(app *gin.RouterGroup, resource *db.Resource) {
	productEntity := repository.NewProductEntity(resource)
	_, _ = productEntity.CreateIndex()

	productRoute := app.Group("product")
	productRoute.GET("", getAll(productEntity))
	productRoute.POST("", middlewares.RequireAuthenticated(), createProduct(productEntity))
	productRoute.GET("/:productId", getProductById(productEntity))
	productRoute.PUT("/:productId", middlewares.RequireAuthenticated(), updateProductById(productEntity))
	productRoute.GET("/serial-number/:serialNumber", getProductBySerialNumber(productEntity))
	productRoute.GET("/:productId/lot", getLotAllByProductId(productEntity))
	productRoute.PUT("/lot/:lotId", updateLotById(productEntity))
	productRoute.GET("/lot/:lotId", getLotById(productEntity))
}

func getAll(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, code, err := productEntity.GetAll()
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func createProduct(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Product{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		device, code, err := productEntity.CreateOne(request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, device)
	}
}

func updateProductById(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("productId")
		request := form.UpdateProduct{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		device, code, err := productEntity.UpdateOneById(id, request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, device)
	}
}

func getProductBySerialNumber(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		serialNumber := ctx.Param("serialNumber")
		result, code, err := productEntity.GetOneBySerialNumber(serialNumber)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func getProductById(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("productId")
		result, code, err := productEntity.GetOneById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func getLotAllByProductId(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productId := ctx.Param("productId")
		result, code, err := productEntity.GetLotAllByProductId(productId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func getLotById(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("lotId")
		result, code, err := productEntity.GetLotOneById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func updateLotById(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("lotId")
		request := form.ProductLot{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		device, code, err := productEntity.UpdateLotOneById(id, request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, device)
	}
}
