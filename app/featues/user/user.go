package user

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/db"
	"mgo-gin/middlewares"
	"mgo-gin/utils/bcrypt"
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
	userRoute.GET("/info", getUserInfo(userEntity))
	userRoute.PUT("/info", updateUserInfo(userEntity))
	userRoute.GET("/keep-alive", keepAlive(userEntity))

	// ADMIN
	userRoute.GET("/:id", middlewares.RequireAuthorization(constant.ADMIN), getUserById(userEntity))
	userRoute.DELETE("/:id", middlewares.RequireAuthorization(constant.ADMIN), deleteUserById(userEntity))
	userRoute.GET("", middlewares.RequireAuthorization(constant.ADMIN), getAllUser(userEntity))
	userRoute.POST("", middlewares.RequireAuthorization(constant.ADMIN), addUser(userEntity))
	userRoute.PUT("/:id", middlewares.RequireAuthorization(constant.ADMIN), updateUser(userEntity))
}

func login(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userRequest := form.Login{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, code, _ := userEntity.GetOneByUsername(userRequest.Username)
		if (user == nil) || bcrypt.ComparePasswordAndHashedPassword(userRequest.Password, user.Password) != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Wrong username or password"})
			return
		}
		token := middlewares.GenerateJwtToken(*user)
		response := gin.H{
			"accessToken": token,
		}
		ctx.JSON(code, response)
	}
}

func signUp(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userRequest := form.User{}
		err := ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, code, err := userEntity.CreateOne(userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func getAllUser(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, _ := ctx.Get("UserId")
		_, code, err := userEntity.GetOneById(userId.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		list, code, err := userEntity.GetAll()
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, list)
	}
}

func getUserInfo(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, _ := ctx.Get("UserId")
		user, code, err := userEntity.GetOneById(userId.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func keepAlive(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, _ := ctx.Get("UserId")
		user, code, err := userEntity.GetOneById(userId.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		token := middlewares.GenerateJwtToken(*user)
		response := gin.H{
			"accessToken": token,
		}
		ctx.JSON(code, response)
	}
}

func getUserById(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, _ := ctx.Get("UserId")
		_, code, err := userEntity.GetOneById(userId.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		user, code, err := userEntity.GetOneById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func deleteUserById(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, _ := ctx.Get("UserId")
		found, code, err := userEntity.GetOneById(userId.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		if found.Id.Hex() == id {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't delete self user"})
			return
		}
		user, code, err := userEntity.RemoveOneById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func addUser(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, _ := ctx.Get("UserId")
		_, code, err := userEntity.GetOneById(userId.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.User{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.CreatedBy = userId.(string)
		user, code, err := userEntity.CreateOne(userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func updateUser(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, _ := ctx.Get("UserId")
		_, code, err := userEntity.GetOneById(userId.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		userRequest := form.User{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId.(string)
		user, code, err := userEntity.UpdateUserById(id, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func updateUserInfo(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, _ := ctx.Get("UserId")
		_, code, err := userEntity.GetOneById(userId.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.User{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId.(string)
		user, code, err := userEntity.UpdateUserById(userId.(string), userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}
