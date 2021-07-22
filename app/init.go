package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mgo-gin/app/featues/notification"
	"mgo-gin/app/featues/order"
	"mgo-gin/app/featues/product"
	"mgo-gin/app/featues/user"
	"mgo-gin/db"
	"mgo-gin/middlewares"
	"net/http"
	"os"
)

type Routes struct {
}

func (app Routes) StartGin() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))

	resource, err := db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	publicRoute := r.Group("/api/v1")

	user.ApplyUserAPI(publicRoute, resource)
	notification.ApplyNotificationAPI(publicRoute, resource)
	product.ApplyProductAPI(publicRoute, resource)
	order.ApplyOrderAPI(publicRoute, resource)

	r.NoRoute(func(context *gin.Context) {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Service Missing / Not found."})
	})

	_ = r.Run(":" + os.Getenv("PORT"))
}
