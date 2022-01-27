package app

import (
	"devper/app/featues/category"
	"devper/app/featues/category/repository"
	"devper/app/featues/notification"
	repository2 "devper/app/featues/notification/repository"
	"devper/app/featues/order"
	repository3 "devper/app/featues/order/repository"
	"devper/app/featues/product"
	repository4 "devper/app/featues/product/repository"
	"devper/app/featues/user"
	repository5 "devper/app/featues/user/repository"
	"devper/db"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

type Routes struct {
}

func (app Routes) StartGin() {
	r := gin.New()

	err := r.SetTrustedProxies(nil)
	if err != nil {
		logrus.Error(err)
	}

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))

	resource, err := db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	publicRoute := r.Group("/api/v1")

	userEntity := repository5.NewUserEntity(resource)
	productEntity := repository4.NewProductEntity(resource)
	orderEntity := repository3.NewOrderEntity(resource)
	notificationEntity := repository2.NewNotificationEntity(resource)
	categoryEntity := repository.NewCategoryEntity(resource)

	user.ApplyAuthAPI(publicRoute, userEntity)
	user.ApplyUserAPI(publicRoute, userEntity)
	notification.ApplyNotificationAPI(publicRoute, notificationEntity)
	product.ApplyProductAPI(publicRoute, productEntity)
	order.ApplyOrderAPI(publicRoute, orderEntity, productEntity)
	category.ApplyCategoryAPI(publicRoute, categoryEntity)

	r.NoRoute(middlewares.NoRoute())

	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		logrus.Error(err)
	}
}
