package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/core/bcrypt"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/middlewares"
	"mgo-gin/utils/constant"
	"net/http"
)

func Login(userEntity repository.IUser) gin.HandlerFunc {
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
