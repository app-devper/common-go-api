package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mgo-gin/app/featues/user/model"
	"net/http"
	"strings"
	"time"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("uit_secret_key")

type Claims struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateJwtToken(user model.User) string {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		UserId: user.Id.Hex(),
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			Audience:  "user",
			Issuer:    "uit",
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
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
		// Initialize a new instance of `Claims`
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(jwtToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if tkn == nil || !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalid authorization header"})
			return
		}
		c.Set("UserId", claims.UserId)
		c.Set("Role", claims.Role)
		logrus.Info("UserId: " + claims.UserId)
		return
	}
}
