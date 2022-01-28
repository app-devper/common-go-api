package usecase

import (
	"devper/app/featues/order/form"
	"devper/app/featues/order/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOrderItemRange(orderEntity repository.IOrder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.GetOrderRange{}
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := orderEntity.GetOrderItemRange(request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
