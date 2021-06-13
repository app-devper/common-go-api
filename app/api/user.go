package api

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/form"
	"mgo-gin/app/repository"
	"mgo-gin/db"
	"mgo-gin/middlewares"
	"mgo-gin/utils/constant"
	"net/http"
)

func ApplyUserAPI(app *gin.RouterGroup, resource *db.Resource) {
	userEntity := repository.NewUserEntity(resource)
	authRoute := app.Group("auth")
	authRoute.POST("/login", login(userEntity))
	authRoute.POST("/sign-up", signUp(userEntity))

	userRoute := app.Group("/user")
	userRoute.Use(middlewares.RequireAuthenticated())
	userRoute.Use(middlewares.RequireAuthorization(constant.ADMIN))
	userRoute.GET("/info", getUserInfo(userEntity))
	userRoute.GET("/:id", getUserById(userEntity))
	userRoute.GET("", getAllUser(userEntity))
}

func login(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		userRequest := form.User{}
		if err := ctx.Bind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, code, _ := userEntity.GetOneByUsername(userRequest.Username)

		if (user == nil) || userRequest.Password != user.Password {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong username or password"})
			return
		}

		token := middlewares.GenerateJwtToken(*user)
		response := gin.H{
			"token": token,
		}
		ctx.JSON(code, response)
	}
}

func signUp(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userRequest := form.User{}
		if err := ctx.Bind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, code, err := userEntity.CreateOne(userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		response := gin.H{
			"user": user,
		}
		ctx.JSON(code, response)
	}
}

func getAllUser(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		list, code, err := userEntity.GetAll()
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		response := gin.H{
			"users": list,
		}
		ctx.JSON(code, response)
	}
}

func getUserInfo(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		username, _ := ctx.Get("Username")
		user, code, err := userEntity.GetOneByUsername(username.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		response := gin.H{
			"user": user,
		}
		ctx.JSON(code, response)
	}
}

func getUserById(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		user, code, err := userEntity.GetOneById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		response := gin.H{
			"user": user,
		}
		ctx.JSON(code, response)
	}
}
