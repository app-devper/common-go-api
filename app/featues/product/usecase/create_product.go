package usecase

import (
	"devper/app/featues/product/form"
	"devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateProduct(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Product{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := productEntity.CreateProduct(request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
