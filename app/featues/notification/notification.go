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

func subscription(notificationEntity repository.INotification) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Subscription{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		device, code, err := notificationEntity.Subscription(request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, device)
	}
}
