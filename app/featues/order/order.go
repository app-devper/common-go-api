package order

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/order/form"
	"mgo-gin/app/featues/order/repository"
	repository2 "mgo-gin/app/featues/product/repository"
	"mgo-gin/db"
	"net/http"
)

func ApplyOrderAPI(app *gin.RouterGroup, resource *db.Resource) {
	orderEntity := repository.NewOrderEntity(resource)
	productEntity := repository2.NewProductEntity(resource)
	orderRoute := app.Group("order")
	orderRoute.POST("", createOrder(orderEntity, productEntity))
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

		ctx.JSON(code, result)
	}
}
