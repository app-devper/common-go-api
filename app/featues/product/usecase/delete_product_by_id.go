package usecase

import (
	"devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteProductById(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("productId")
		result, err := productEntity.RemoveProductById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
