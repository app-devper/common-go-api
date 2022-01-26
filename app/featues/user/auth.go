package user

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/app/featues/user/usecase"
	"mgo-gin/middlewares"
)

func ApplyAuthAPI(app *gin.RouterGroup, userEntity repository.IUser) {

	authRoute := app.Group("auth")

	authRoute.POST("/login",
		usecase.Login(userEntity),
	)

	authRoute.POST("/sign-up",
		usecase.SignUp(userEntity),
	)

	authRoute.POST("/verification/user",
		usecase.VerifyUser(userEntity),
	)

	authRoute.POST("/verification/request",
		usecase.VerifyRequest(userEntity),
	)

	authRoute.POST("/verification/code",
		usecase.VerifyCode(userEntity),
	)

	authRoute.GET("/verification/info",
		middlewares.RequireActionToken(),
		usecase.GetUserRefInfo(userEntity),
	)

	authRoute.POST("/set-password",
		middlewares.RequireActionToken(),
		usecase.SetPassword(userEntity),
	)

}
