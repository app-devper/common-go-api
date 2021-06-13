package middlewares

import (
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRecovery() gin.HandlerFunc {
	return nice.Recovery(recoveryHandler)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": err,
	})
}
