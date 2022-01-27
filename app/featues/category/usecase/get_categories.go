package usecase

import (
	"devper/app/featues/category/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetCategories(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := entity.GetCategoryAll()
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
