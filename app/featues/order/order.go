package order

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/order/form"
	"mgo-gin/app/featues/order/repository"
	"mgo-gin/db"
	"net/http"
)

func ApplyOrderAPI(app *gin.RouterGroup, resource *db.Resource) {
	orderEntity := repository.NewOrderEntity(resource)

	orderRoute := app.Group("order")
	orderRoute.POST("", createOrder(orderEntity))
}

func createOrder(orderEntity repository.IOrder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Order{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		device, code, err := orderEntity.CreateOrder(request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, device)
	}
}
