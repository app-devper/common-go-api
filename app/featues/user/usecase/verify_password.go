package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/core/bcrypt"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/utils/constant"
	"net/http"
)

func VerifyPassword(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		user, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.VerifyPassword{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if (user == nil) || bcrypt.ComparePasswordAndHashedPassword(userRequest.Password, user.Password) != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
			return
		}
		if user.Status != constant.ACTIVE {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not active"})
			return
		}
		response := gin.H{
			"message": "success",
		}
		ctx.JSON(code, response)
	}
}
