package order

import (
	"devper/app/core/constant"
	"devper/app/featues/order/repository"
	"devper/app/featues/order/usecase"
	repository2 "devper/app/featues/product/repository"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
)

func ApplyOrderAPI(app *gin.RouterGroup, orderEntity repository.IOrder, productEntity repository2.IProduct) {
	orderRoute := app.Group("order")

	orderRoute.POST("",
		usecase.CreateOrder(orderEntity, productEntity),
	)

	orderRoute.GET("",
		usecase.GetOrdersRange(orderEntity),
	)

	orderRoute.GET("/:orderId",
		usecase.GetOrderById(orderEntity),
	)

	orderRoute.DELETE("/:orderId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.DeleteOrderById(orderEntity, productEntity),
	)

	orderRoute.GET("/:orderId/total-cost",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.UpdateTotalCostById(orderEntity, productEntity),
	)

	orderRoute.GET("/item",
		usecase.GetOrderItemRange(orderEntity),
	)

	orderRoute.GET("/item/:itemId",
		usecase.GetOrderItemById(orderEntity),
	)

	orderRoute.DELETE("/item/:itemId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.DeleteOrderItemById(orderEntity, productEntity),
	)

	orderRoute.GET("/product/:productId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.GetOrderItemByProductId(orderEntity),
	)

	orderRoute.DELETE("/:orderId/product/:productId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.DeleteOrderItemByOrderProductId(orderEntity, productEntity),
	)

}
