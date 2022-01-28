package middlewares

import (
	"devper/app/core/constant"
	"devper/app/featues/user/repository"
	"errors"
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

type AccessClaims struct {
	UserRefId string `json:"userRefId"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

type ActionClaims struct {
	UserRefId string `json:"userRefId"`
	Objective string `json:"objective"`
	jwt.StandardClaims
}

func GenerateJwtToken(userRefId string, role string, expirationTime time.Time) string {
	claims := &AccessClaims{
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

func RequireAuthenticated(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			err := errors.New("missing authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		jwtToken := strings.Split(token, "Bearer ")
		if len(jwtToken) < 2 {
			err := errors.New("missing authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		claims := &AccessClaims{}
		tkn, err := jwt.ParseWithClaims(jwtToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if tkn == nil || !tkn.Valid || claims.UserRefId == "" {
			err := errors.New("token invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		userRef, err := userEntity.GetVerificationById(claims.UserRefId)
		if userRef == nil {
			err := errors.New("user ref invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if userRef.Status != constant.ACTIVE {
			err := errors.New("user ref not active")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if userRef.Objective != constant.AccessApi {
			err := errors.New("objective invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if userRef.ExpireDate.Before(time.Now()) {
			err := errors.New("token invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		user, err := userEntity.GetUserById(userRef.UserId.Hex())
		if err != nil {
			logrus.Error(err)
		}
		if user == nil {
			err := errors.New("user invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if user.Status != constant.ACTIVE {
			err := errors.New("user not active")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("UserRefId", claims.UserRefId)
		ctx.Set("UserId", userRef.UserId.Hex())
		ctx.Set("Role", claims.Role)

		logrus.Info("UserRefId: " + claims.UserRefId)
		logrus.Info("UserId: " + userRef.UserId.Hex())
		logrus.Info("Role: " + claims.Role)
		return
	}
}

func GenerateActionToken(userRefId string, objective string, expirationTime time.Time) string {
	claims := &ActionClaims{
		UserRefId: userRefId,
		Objective: objective,
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

func RequireActionToken(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("X-Action-Token")
		if token == "" {
			err := errors.New("missing action token header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		claims := &ActionClaims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if tkn == nil || !tkn.Valid {
			err := errors.New("token invalid action token header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		userRef, err := userEntity.GetVerificationById(claims.UserRefId)
		if userRef == nil {
			err := errors.New("user ref invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if userRef.Status != constant.ACTIVE {
			err := errors.New("user ref not active")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if userRef.Objective != claims.Objective {
			err := errors.New("objective invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if userRef.ExpireDate.Before(time.Now()) {
			err := errors.New("token invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		user, err := userEntity.GetUserById(userRef.UserId.Hex())
		if user == nil {
			err := errors.New("user invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if user.Status != constant.ACTIVE {
			err := errors.New("user not active")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("UserId", user.Id.Hex())
		ctx.Set("UserRefId", claims.UserRefId)

		logrus.Info("UserRefId: " + claims.UserRefId)
		logrus.Info("UserId: " + user.Id.Hex())
		return
	}
}
