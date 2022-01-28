package usecase

import (
	"devper/app/featues/category/form"
	"devper/app/featues/category/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateCategoryById(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		categoryId := ctx.Param("categoryId")
		request := form.Category{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := entity.UpdateCategoryById(categoryId, request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
