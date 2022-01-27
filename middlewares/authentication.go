package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

// Create the JWT key used to create the signature
var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	UserRefId string `json:"userRefId"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

type ActionClaims struct {
	VerifyId string `json:"verifyId"`
	jwt.StandardClaims
}

func GenerateJwtToken(userRefId string, role string, expirationTime time.Time) string {
	claims := &Claims{
		UserRefId: userRefId,
		Role:      role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Audience:  "user",
			Issuer:    "uit",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logrus.Error(err)
	}
	return tokenString
}

func RequireAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			return
		}
		jwtToken := strings.Split(token, "Bearer ")
		if len(jwtToken) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(jwtToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			logrus.Error(err)
		}
		if tkn == nil || !tkn.Valid || claims.UserRefId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalid authorization header"})
			return
		}
		c.Set("UserRefId", claims.UserRefId)
		c.Set("Role", claims.Role)
		logrus.Info("UserRefId: " + claims.UserRefId)
		return
	}
}

func GenerateActionToken(verifyId string, expirationTime time.Time) string {
	claims := &ActionClaims{
		VerifyId: verifyId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logrus.Error(err)
	}
	return tokenString
}

func RequireActionToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Action-Token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing action token header"})
			return
		}
		claims := &ActionClaims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			logrus.Error(err)
		}
		if tkn == nil || !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalid action token header"})
			return
		}
		c.Set("verifyId", claims.VerifyId)
		logrus.Info("VerifyId: " + claims.VerifyId)
		return
	}
}
