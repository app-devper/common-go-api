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

	authRoute.POST("/verify-user",
		usecase.VerifyUser(userEntity),
	)

	authRoute.POST("/verify-request",
		usecase.VerifyRequest(userEntity),
	)

	authRoute.POST("/verify-code",
		usecase.VerifyCode(userEntity),
	)

	authRoute.GET("/verify-info",
		middlewares.RequireActionToken(),
		usecase.GetUserRefInfo(userEntity),
	)

	authRoute.POST("/set-password",
		middlewares.RequireActionToken(),
		usecase.SetPassword(userEntity),
	)

}
