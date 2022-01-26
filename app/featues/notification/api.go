package notification

import (
	"devper/app/featues/notification/repository"
	"devper/app/featues/notification/usecase"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
)

func ApplyNotificationAPI(app *gin.RouterGroup, notificationEntity repository.INotification) {
	notificationRoute := app.Group("notification")

	notificationRoute.POST("/subscription",
		middlewares.RequireAuthenticated(),
		usecase.Subscription(notificationEntity),
	)

	notificationRoute.POST("/notify",
		usecase.NotifyMessage(),
	)
}
