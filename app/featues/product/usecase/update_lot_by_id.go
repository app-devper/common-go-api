package usecase

import (
	"devper/app/featues/product/form"
	"devper/app/featues/product/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateLotById(productEntity repository.IProduct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("lotId")
		request := form.ProductLot{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := productEntity.UpdateLotById(id, request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
