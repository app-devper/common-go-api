package usecase

import (
	"devper/app/core"
	"devper/app/featues/order/form"
	"devper/app/featues/order/repository"
	repository2 "devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func CreateOrder(orderEntity repository.IOrder, productEntity repository2.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Order{}
		if err := ctx.ShouldBind(&request); err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		totalCost := 0.0
		for index, item := range request.Items {
			request.Items[index].CostPrice = productEntity.GetTotalCostPrice(item.ProductId, item.Quantity)
			totalCost += request.Items[index].CostPrice
		}
		request.TotalCost = totalCost

		result, err := orderEntity.CreateOrder(request)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, item := range request.Items {
			_, _ = productEntity.RemoveQuantityById(item.ProductId, item.Quantity)
		}

		if request.Message != "" {
			date := core.ToFormat(result.CreatedDate)
			_, _ = core.NotifyMassage("รายการวันที่ " + date + "\n\n" + request.Message)
		}

		ctx.JSON(http.StatusOK, result)
	}
}
