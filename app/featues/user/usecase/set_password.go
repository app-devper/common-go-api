package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"net/http"
)

func SetPassword(userEntity repository.IUser) gin.HandlerFunc {
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
