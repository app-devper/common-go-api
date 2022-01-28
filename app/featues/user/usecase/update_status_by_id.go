package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateStatusById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.UpdateStatus{}
		err := ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		userId := ctx.GetString("UserId")
		if userId == id {
			err = errors.New("can't update self user")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		result, err := userEntity.UpdateStatusById(id, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
