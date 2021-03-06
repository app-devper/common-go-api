package usecase

import (
	"devper/app/featues/category/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCategoryById(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		categoryId := ctx.Param("categoryId")
		result, err := entity.GetCategoryById(categoryId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
