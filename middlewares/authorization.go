package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetRolesFromToken(tokenReq string) (role []string) {
	token, _ := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	claim := token.Claims.(jwt.MapClaims)
	var roles []string
	rolesResource := claim["role"].(interface{})
	roles = append(roles, rolesResource.(string))
	if len(roles) <= 0 {
		return nil
	}
	return roles
}

func RequireAuthorization(auths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		jwtToken := strings.Split(token, "Bearer ")
		roles := GetRolesFromToken(jwtToken[1])
		if len(roles) <= 0 {
			invalidRequest(c)
			return
		}
		isAccessible := false
		if len(roles) < len(auths) || len(roles) == len(auths) {
			for _, auth := range auths {
				for _, role := range roles {
					if role == auth {
						isAccessible = true
						break
					}
				}
			}
		}
		if len(roles) > len(auths) {
			for _, role := range roles {
				for _, auth := range auths {
					if auth == role {
						isAccessible = true
						break
					}
				}
			}
		}
		if isAccessible == false {
			notPermission(c)
			return
		}
		c.Next()
	}
}

func invalidRequest(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid request, restricted endpoint"})
}

func notPermission(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Dont have permission"})
}
