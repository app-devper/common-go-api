package usecase

import (
	"devper/app/core/utils"
	"devper/app/featues/order/repository"
	repository2 "devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteOrderItemById(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		itemId := ctx.Param("itemId")
		result, err := orderEntity.RemoveOrderItemById(itemId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = orderEntity.UpdateTotalOrderById(result.OrderId.Hex())
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, _ = productEntity.AddQuantityById(result.ProductId.Hex(), result.Quantity)

		date := utils.ToFormat(result.CreatedDate)
		_, _ = utils.NotifyMassage("ยกเลิกสินค้ารายการวันที่ " + date + "\n\n1. " + result.GetMessage())

		ctx.JSON(http.StatusOK, result)
	}
}
