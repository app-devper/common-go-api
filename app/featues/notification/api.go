package notification

import (
	"devper/app/featues/notification/repository"
	"devper/app/featues/notification/usecase"
	repository2 "devper/app/featues/user/repository"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
)

func ApplyNotificationAPI(
	app *gin.RouterGroup,
	notificationEntity repository.INotification,
	userEntity repository2.IUser,
) {
	notificationRoute := app.Group("notification")

	notificationRoute.POST("/subscription",
		middlewares.RequireAuthenticated(userEntity),
		usecase.Subscription(notificationEntity),
	)

	notificationRoute.POST("/notify",
		usecase.NotifyMessage(),
	)
}
