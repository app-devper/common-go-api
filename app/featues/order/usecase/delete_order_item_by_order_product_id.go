package usecase

import (
	"devper/app/core"
	"devper/app/featues/order/repository"
	repository2 "devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func DeleteOrderItemByOrderProductId(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId := ctx.Param("orderId")
		productId := ctx.Param("productId")
		result, err := orderEntity.RemoveOrderItemByOrderProductId(orderId, productId)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = orderEntity.UpdateTotalOrderById(orderId)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, _ = productEntity.AddQuantityById(productId, result.Quantity)

		date := core.ToFormat(result.CreatedDate)
		_, _ = core.NotifyMassage("ยกเลิกสินค้ารายการวันที่ " + date + "\n\n1. " + result.GetMessage())

		ctx.JSON(http.StatusOK, result)
	}
}
