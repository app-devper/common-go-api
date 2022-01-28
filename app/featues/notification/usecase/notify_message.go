package usecase

import (
	"devper/app/core/utils"
	"devper/app/featues/notification/form"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotifyMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Notify{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := utils.NotifyMassage(request.Message)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
