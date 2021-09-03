package order

import (
	"fmt"
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
	orderRoute.GET("total", updateTotal(orderEntity))
	orderRoute.GET("/:id", getOrderById(orderEntity))
	orderRoute.GET("/:id/total-cost", updateTotalCost(orderEntity, productEntity))
	orderRoute.GET("/product/:productId", middlewares.RequireAuthenticated(), getOrderItemByProductId(orderEntity))
	orderRoute.DELETE("/:id/product/:productId", middlewares.RequireAuthenticated(), deleteProductInOrder(orderEntity, productEntity))
	orderRoute.DELETE("/:id", middlewares.RequireAuthenticated(), deleteOrder(orderEntity, productEntity))
}

func createOrder(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Order{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for index, item := range request.Items {
			request.Items[index].CostPrice = productEntity.GetTotalCostPrice(item.ProductId, item.Quantity)
			request.TotalCost += item.CostPrice
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

func updateTotal(orderEntity repository.IOrder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, code, err := orderEntity.UpdateTotal()
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func updateTotalCost(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		order, code, err := orderEntity.GetOrderDetailById(id)
		totalCost := 0.0
		for _, item := range order.Items {
			orderItem := form.OrderItem{
				CostPrice: productEntity.GetTotalCostPrice(item.ProductId.Hex(), item.Quantity),
				Quantity:  item.Quantity,
				Price:     item.Price,
				Discount:  item.Discount,
				ProductId: item.ProductId.Hex(),
			}
			_, _, _ = orderEntity.UpdateOrderItemById(item.Id.Hex(), orderItem)
			totalCost += orderItem.CostPrice
		}
		result, code, err := orderEntity.UpdateTotalCostOrderById(id, totalCost)
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
		result, code, err := orderEntity.GetOrderDetailById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func getOrderItemByProductId(orderEntity repository.IOrder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productId := ctx.Param("productId")
		result, code, err := orderEntity.GetOrderItemByProductId(productId)
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

		var message = ""
		var no = 1
		for _, item := range result.Items {
			_, _, _ = productEntity.AddQuantityById(item.ProductId.Hex(), item.Quantity)
			message += fmt.Sprintf("%d. %s\n", no, item.GetMessage())
			no += 1
		}
		message += fmt.Sprintf("\nรวม %.2f บาท", result.Total)

		location, _ := time.LoadLocation("Asia/Bangkok")
		format := "02 Jan 2006 15:04"
		date := result.CreatedDate.In(location).Format(format)
		_, _ = core.NotifyMassage("ยกเลิกรายการวันที่ " + date + "\n\n" + message)

		ctx.JSON(code, result)
	}
}

func deleteProductInOrder(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		productId := ctx.Param("productId")
		result, code, err := orderEntity.RemoveProductByOrderProductId(id, productId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}

		_, _, _ = productEntity.AddQuantityById(productId, result.Quantity)

		location, _ := time.LoadLocation("Asia/Bangkok")
		format := "02 Jan 2006 15:04"
		date := result.CreatedDate.In(location).Format(format)
		_, _ = core.NotifyMassage("ยกเลิกสินค้ารายการวันที่ " + date + "\n\n1. " + result.GetMessage())

		ctx.JSON(code, result)
	}
}
