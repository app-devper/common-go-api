package usecase

import (
	"devper/app/core/utils"
	"devper/app/featues/order/repository"
	repository2 "devper/app/featues/product/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteOrderById(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId := ctx.Param("orderId")
		result, err := orderEntity.RemoveOrderById(orderId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var message = ""
		var no = 1
		for _, item := range result.Items {
			_, _ = productEntity.AddQuantityById(item.ProductId.Hex(), item.Quantity)
			message += fmt.Sprintf("%d. %s\n", no, item.GetMessage())
			no += 1
		}
		message += fmt.Sprintf("\nรวม %.2f บาท", result.Total)

		date := utils.ToFormat(result.CreatedDate)
		_, _ = utils.NotifyMassage("ยกเลิกรายการวันที่ " + date + "\n\n" + message)

		ctx.JSON(http.StatusOK, result)
	}
}
