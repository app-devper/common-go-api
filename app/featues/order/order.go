package order

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/core"
	"mgo-gin/app/featues/order/form"
	"mgo-gin/app/featues/order/repository"
	repository2 "mgo-gin/app/featues/product/repository"
	"mgo-gin/db"
	"mgo-gin/middlewares"
	"net/http"
	"time"
)

func ApplyOrderAPI(app *gin.RouterGroup, resource *db.Resource) {
	orderEntity := repository.NewOrderEntity(resource)
	productEntity := repository2.NewProductEntity(resource)
	orderRoute := app.Group("order")
	orderRoute.POST("", createOrder(orderEntity, productEntity))
	orderRoute.GET("", getOrderRange(orderEntity))
	orderRoute.GET("/:id", getOrderById(orderEntity))
	orderRoute.DELETE("/:id", middlewares.RequireAuthenticated(), deleteOrder(orderEntity, productEntity))
}

func createOrder(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Order{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, code, err := orderEntity.CreateOrder(request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		for _, item := range request.Items {
			_, _, _ = productEntity.RemoveQuantityById(item.ProductId, item.Quantity)
		}
		if request.Message != "" {
			location, _ := time.LoadLocation("Asia/Bangkok")
			format := "02 Jan 2006 15:04"
			date := result.CreatedDate.In(location).Format(format)
			_, _ = core.NotifyMassage("รายการวันที่ " + date + "\n\n" + request.Message)
		}
		ctx.JSON(code, result)
	}
}

func getOrderRange(orderEntity repository.IOrder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.GetOrder{}
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, code, err := orderEntity.GetOrderRange(request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func getOrderById(orderEntity repository.IOrder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		result, code, err := orderEntity.GetOrderById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func deleteOrder(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		result, code, err := orderEntity.RemoveOrderById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		for _, item := range result.Items {
			_, _, _ = productEntity.AddQuantityById(item.ProductId.Hex(), item.Quantity)
		}
		if result.Message != "" {
			location, _ := time.LoadLocation("Asia/Bangkok")
			format := "02 Jan 2006 15:04"
			date := result.CreatedDate.In(location).Format(format)
			_, _ = core.NotifyMassage("ยกเลิกรายการวันที่ " + date + "\n\n" + result.Message)
		}
		ctx.JSON(code, result)
	}
}
