package usecase

import (
	"devper/app/featues/category/form"
	"devper/app/featues/category/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateCategory(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Category{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := entity.CreateCategory(request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
