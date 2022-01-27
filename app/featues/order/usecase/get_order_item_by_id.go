package usecase

import (
	"devper/app/featues/order/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetOrderItemById(orderEntity repository.IOrder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		itemId := ctx.Param("itemId")
		result, err := orderEntity.GetOrderItemDetailById(itemId)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
