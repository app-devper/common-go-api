package usecase

import (
	"devper/app/featues/notification/form"
	"devper/app/featues/notification/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Subscription(notificationEntity repository.INotification) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Subscription{}
		if err := ctx.ShouldBind(&request); err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := notificationEntity.Subscription(request)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
