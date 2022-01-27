package usecase

import (
	"devper/app/featues/order/form"
	"devper/app/featues/order/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetOrdersRange(orderEntity repository.IOrder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.GetOrderRange{}
		if err := ctx.ShouldBindQuery(&request); err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := orderEntity.GetOrderRange(request)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
