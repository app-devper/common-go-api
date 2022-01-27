package usecase

import (
	"devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetLotsByProductId(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productId := ctx.Param("productId")
		result, err := productEntity.GetLotAllByProductId(productId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
