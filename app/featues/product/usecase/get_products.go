package usecase

import (
	"devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProducts(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := productEntity.GetProductAll()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
