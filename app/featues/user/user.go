package user

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/app/featues/user/usecase"
	"mgo-gin/middlewares"
	"mgo-gin/utils/constant"
)

func ApplyUserAPI(app *gin.RouterGroup, userEntity repository.IUser) {

	userRoute := app.Group("/user")

	userRoute.GET("/info",
		middlewares.RequireAuthenticated(),
		usecase.GetUserInfo(userEntity),
	)

	userRoute.PUT("/info",
		middlewares.RequireAuthenticated(),
		usecase.UpdateUserInfo(userEntity),
	)

	userRoute.PUT("/change-password",
		middlewares.RequireAuthenticated(),
		usecase.ChangePassword(userEntity),
	)

	userRoute.GET("/keep-alive",
		middlewares.RequireAuthenticated(),
		usecase.KeepAlive(userEntity),
	)

	userRoute.GET("/verify-password",
		middlewares.RequireAuthenticated(),
		usecase.VerifyPassword(userEntity),
	)

	// ADMIN
	userRoute.GET("",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.GetUsers(userEntity),
	)

	userRoute.POST("",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.AddUser(userEntity),
	)

	userRoute.GET("/:id",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.GetUserById(userEntity),
	)

	userRoute.DELETE("/:id",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.DeleteUserById(userEntity),
	)

	userRoute.PUT("/:id",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.UpdateUserById(userEntity),
	)

	userRoute.PATCH("/:id/status",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.UpdateStatusById(userEntity),
	)

	userRoute.PATCH("/:id/role",
		middlewares.RequireAuthenticated(),
		middlewares.RequireAuthorization(constant.ADMIN),
		usecase.UpdateRoleById(userEntity),
	)
}
