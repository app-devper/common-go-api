package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SetPassword(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		userRef, err := userEntity.GetVerificationById(userRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userRef.Objective != constant.SetPassword {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong objective"})
			return
		}
		if userRef.ExpireDate.Before(time.Now()) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "token invalid"})
			return
		}
		user, err := userEntity.GetUserById(userRef.UserId.Hex())
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.SetPassword{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := userEntity.SetPassword(user.Id.Hex(), userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, _ = userEntity.RevokeVerification(userRefId)
		ctx.JSON(http.StatusOK, result)
	}
}
