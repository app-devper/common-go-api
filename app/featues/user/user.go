package user

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/core/bcrypt"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/db"
	"mgo-gin/middlewares"
	"mgo-gin/utils/constant"
	"net/http"
)

func ApplyUserAPI(app *gin.RouterGroup, resource *db.Resource) {
	userEntity := repository.NewUserEntity(resource)
	_, _ = userEntity.CreateIndex()

	authRoute := app.Group("auth")
	authRoute.POST("/login", login(userEntity))
	authRoute.POST("/sign-up", signUp(userEntity))
	authRoute.POST("/verification/user", verifyUser(userEntity))
	authRoute.POST("/verification/request", verifyRequest(userEntity))
	authRoute.POST("/verification/code", verifyCode(userEntity))
	authRoute.GET("/verification/info", middlewares.RequireActionToken(), getVerifyInfo(userEntity))

	userRoute := app.Group("/user")
	userRoute.GET("/info", middlewares.RequireAuthenticated(), getUserInfo(userEntity))
	userRoute.PUT("/info", middlewares.RequireAuthenticated(), updateUserInfo(userEntity))
	userRoute.PUT("/change-password", middlewares.RequireAuthenticated(), changePassword(userEntity))
	userRoute.GET("/keep-alive", middlewares.RequireAuthenticated(), keepAlive(userEntity))
	userRoute.POST("/set-password", middlewares.RequireActionToken(), setPassword(userEntity))

	// ADMIN
	userRoute.GET("/:id", middlewares.RequireAuthenticated(), middlewares.RequireAuthorization(constant.ADMIN), getUserById(userEntity))
	userRoute.DELETE("/:id", middlewares.RequireAuthenticated(), middlewares.RequireAuthorization(constant.ADMIN), deleteUserById(userEntity))
	userRoute.GET("", middlewares.RequireAuthenticated(), middlewares.RequireAuthorization(constant.ADMIN), getAllUser(userEntity))
	userRoute.POST("", middlewares.RequireAuthenticated(), middlewares.RequireAuthorization(constant.ADMIN), addUser(userEntity))
	userRoute.PUT("/:id", middlewares.RequireAuthenticated(), middlewares.RequireAuthorization(constant.ADMIN), updateUser(userEntity))
	userRoute.PATCH("/:id/status", middlewares.RequireAuthenticated(), middlewares.RequireAuthorization(constant.ADMIN), updateStatusById(userEntity))
	userRoute.PATCH("/:id/role", middlewares.RequireAuthenticated(), middlewares.RequireAuthorization(constant.ADMIN), updateRoleById(userEntity))
}

func login(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.Login{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, code, _ := userEntity.GetUserByUsername(userRequest.Username)
		if (user == nil) || bcrypt.ComparePasswordAndHashedPassword(userRequest.Password, user.Password) != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong username or password"})
			return
		}
		if user.Status != constant.ACTIVE {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not active"})
			return
		}
		token := middlewares.GenerateJwtToken(*user)
		response := gin.H{
			"accessToken": token,
		}
		ctx.JSON(code, response)
	}
}

func signUp(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.User{}
		err := ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, code, err := userEntity.CreateUser(userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func verifyUser(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyUser{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, _, err := userEntity.GetUserByUsername(userRequest.Username)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userRef, code, err := userEntity.CreateVerification(user.Id.Hex(), userRequest.Objective)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		channels := []form.Channel{{
			Channel:     "MOBILE",
			ChannelInfo: user.Phone,
		}, {
			Channel:     "EMAIL",
			ChannelInfo: user.Email,
		}}
		response := gin.H{
			"userRefId":      userRef.Id,
			"verifyChannels": channels,
		}
		ctx.JSON(code, response)
	}
}

func verifyRequest(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyRequest{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRef, code, err := userEntity.GetVerificationById(userRequest.UserRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userRef.Status == constant.ACTIVE {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ref is active"})
			return
		}
		userRef, code, err = userEntity.UpdateVerification(userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, userRef)
	}
}

func verifyCode(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyCode{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRef, code, err := userEntity.GetVerificationById(userRequest.UserRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userRef.Status == constant.ACTIVE {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ref is active"})
			return
		}
		if userRequest.RefId != userRef.RefId || userRequest.Code != userRef.Code {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong code"})
			return
		}
		userRef, code, err = userEntity.ActiveVerification(userRequest.UserRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token := middlewares.GenerateActionToken(userRequest.UserRefId)
		response := gin.H{
			"actionToken": token,
		}
		ctx.JSON(code, response)
	}
}

func getVerifyInfo(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		userRef, code, err := userEntity.GetVerificationById(userRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, code, err := userEntity.GetUserById(userRef.UserId.Hex())
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func setPassword(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		userRef, code, err := userEntity.GetVerificationById(userRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userRef.Objective != "SET_PASSWORD" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong objective"})
			return
		}
		user, code, err := userEntity.GetUserById(userRef.UserId.Hex())
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.SetPassword{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, code, err = userEntity.SetPassword(user.Id.Hex(), userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		_, _, _ = userEntity.RemoveVerification(userRefId)
		ctx.JSON(code, user)
	}
}

func getAllUser(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		list, code, err := userEntity.GetUserAll()
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, list)
	}
}

func getUserInfo(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		user, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func keepAlive(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		user, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		if user.Status != constant.ACTIVE {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not active"})
			return
		}
		token := middlewares.GenerateJwtToken(*user)
		response := gin.H{
			"accessToken": token,
		}
		ctx.JSON(code, response)
	}
}

func getUserById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		user, code, err := userEntity.GetUserById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func deleteUserById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		if userId == id {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't delete self user"})
			return
		}
		user, code, err := userEntity.RemoveUserById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func addUser(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
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
		userRequest.CreatedBy = userId
		user, code, err := userEntity.CreateUser(userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func updateUser(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		userRequest := form.UpdateUser{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		user, code, err := userEntity.UpdateUserById(id, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func updateUserInfo(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.UpdateUser{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		user, code, err := userEntity.UpdateUserById(userId, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func updateStatusById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		if userId == id {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't update self user"})
			return
		}
		userRequest := form.UpdateStatus{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		user, code, err := userEntity.UpdateStatusById(id, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func updateRoleById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		if userId == id {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't update self user"})
			return
		}
		userRequest := form.UpdateRole{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		user, code, err := userEntity.UpdateRoleById(id, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}

func changePassword(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		user, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.ChangePassword{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if (user == nil) || bcrypt.ComparePasswordAndHashedPassword(userRequest.OldPassword, user.Password) != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong password"})
			return
		}
		user, code, err = userEntity.ChangePassword(userId, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}
