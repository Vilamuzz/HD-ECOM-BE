package middleware

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *appMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestToken := c.Request.Header.Get("Authorization")
		if requestToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Missing token", nil, nil))
			return
		}

		if strings.HasPrefix(strings.ToLower(requestToken), "bearer ") {
			requestToken = strings.TrimSpace(requestToken[7:])
		}

		token, err := jwt.Parse(requestToken, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid token", nil, nil))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid token claims", nil, nil))
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Token missing exp", nil, nil))
			return
		}
		if float64(time.Now().Unix()) > exp {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Token expired", nil, nil))
			return
		}

		var user models.User
		c.Set("currentUser", user)
		c.Next()
	}
}
