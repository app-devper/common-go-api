package usecase

import (
	"devper/app/featues/order/form"
	"devper/app/featues/order/repository"
	repository2 "devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func UpdateTotalCostById(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId := ctx.Param("orderId")
		order, err := orderEntity.GetOrderDetailById(orderId)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		totalCost := 0.0
		for _, item := range order.Items {
			orderItem := form.OrderItem{
				CostPrice: productEntity.GetTotalCostPrice(item.ProductId.Hex(), item.Quantity),
				Quantity:  item.Quantity,
				Price:     item.Price,
				Discount:  item.Discount,
				ProductId: item.ProductId.Hex(),
			}
			_, _ = orderEntity.UpdateOrderItemById(item.Id.Hex(), orderItem)
			totalCost += orderItem.CostPrice
		}
		result, err := orderEntity.UpdateTotalCostOrderById(orderId, totalCost)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
