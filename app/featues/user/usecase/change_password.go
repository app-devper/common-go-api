package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/core/bcrypt"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"net/http"
)

func ChangePassword(userEntity repository.IUser) gin.HandlerFunc {
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
