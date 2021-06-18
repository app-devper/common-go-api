package notification

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/notification/form"
	"mgo-gin/app/featues/notification/repository"
	"mgo-gin/db"
	"mgo-gin/middlewares"
	"net/http"
)

func ApplyNotificationAPI(app *gin.RouterGroup, resource *db.Resource) {
	notificationEntity := repository.NewNotificationEntity(resource)
	notificationRoute := app.Group("notification")
	notificationRoute.Use(middlewares.RequireAuthenticated())
	notificationRoute.POST("/subscription", subscription(notificationEntity))
}

func subscription(notificationEntity repository.INotification) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		request := form.Subscription{}
		if err := ctx.Bind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		subscription, code, err := notificationEntity.Subscription(request)
		if err != nil {
			ctx.JSON(code, gin.H{"error": err.Error()})
			return
		}
		response := gin.H{
			"subscription": subscription,
		}
		ctx.JSON(code, response)
	}
}
