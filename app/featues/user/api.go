package user

import (
	"devper/app/core/constant"
	"devper/app/featues/user/repository"
	"devper/app/featues/user/usecase"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
)

func ApplyAuthAPI(
	app *gin.RouterGroup,
	userEntity repository.IUser,
) {

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

	authRoute.POST("/verify-channel",
		usecase.VerifyUserChannel(userEntity),
	)

	authRoute.POST("/verify-code",
		usecase.VerifyUserCode(userEntity),
	)

	authRoute.GET("/verify-info",
		middlewares.RequireActionToken(userEntity),
		usecase.GetVerifyInfo(userEntity),
	)

	authRoute.POST("/verify-password",
		middlewares.RequireAuthenticated(userEntity),
		usecase.VerifyPassword(userEntity),
	)

	authRoute.POST("/logout",
		middlewares.RequireAuthenticated(userEntity),
		usecase.Logout(userEntity),
	)

}

func ApplyUserAPI(
	app *gin.RouterGroup,
	userEntity repository.IUser,
) {

	userRoute := app.Group("/user")

	userRoute.GET("/info",
		middlewares.RequireAuthenticated(userEntity),
		usecase.GetUserInfo(userEntity),
	)

	userRoute.PUT("/info",
		middlewares.RequireAuthenticated(userEntity),
		usecase.UpdateUserInfo(userEntity),
	)

	userRoute.PUT("/change-password",
		middlewares.RequireAuthenticated(userEntity),
		usecase.ChangePassword(userEntity),
	)

	userRoute.GET("/keep-alive",
		middlewares.RequireAuthenticated(userEntity),
		usecase.KeepAlive(userEntity),
	)

	userRoute.POST("/set-password",
		middlewares.RequireActionToken(userEntity),
		usecase.SetPassword(userEntity),
	)

	// ADMIN
	userRoute.GET("",
		middlewares.RequireAuthenticated(userEntity),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.GetUsers(userEntity),
	)

	userRoute.POST("",
		middlewares.RequireAuthenticated(userEntity),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.AddUser(userEntity),
	)

	userRoute.GET("/:id",
		middlewares.RequireAuthenticated(userEntity),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.GetUserById(userEntity),
	)

	userRoute.DELETE("/:id",
		middlewares.RequireAuthenticated(userEntity),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.DeleteUserById(userEntity),
	)

	userRoute.PUT("/:id",
		middlewares.RequireAuthenticated(userEntity),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.UpdateUserById(userEntity),
	)

	userRoute.PATCH("/:id/status",
		middlewares.RequireAuthenticated(userEntity),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.UpdateStatusById(userEntity),
	)

	userRoute.PATCH("/:id/role",
		middlewares.RequireAuthenticated(userEntity),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.UpdateRoleById(userEntity),
	)
}
