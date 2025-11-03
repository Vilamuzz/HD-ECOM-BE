package middleware

import (
	"app/domain/models"
	"app/helpers"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func (m *appMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from Authorization header first
		requestToken := c.Request.Header.Get("Authorization")

		// If not in header, try query parameter (for WebSocket)
		if requestToken == "" {
			requestToken = c.Query("token")
		} else {
			log.Printf("[AUTH] Token from Authorization header: %s", requestToken)
		}

		if requestToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Missing token", nil, nil))
			return
		}

		// Remove "Bearer " prefix if present
		if strings.HasPrefix(strings.ToLower(requestToken), "bearer ") {
			requestToken = strings.TrimSpace(requestToken[7:])
		}

		// Parse and validate JWT
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

		// Check expiration
		exp, _ := claims["exp"].(float64)
		if float64(time.Now().Unix()) > exp {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Token expired", nil, nil))
			return
		}

		userID, err := strconv.ParseInt(claims["user_id"].(string), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid user_id format", nil, nil))
			return
		}

		// Extract other user info
		username, _ := claims["username"].(string)
		email, _ := claims["email"].(string)

		// Handle role as either string or number
		var role string
		if r, ok := claims["role"].(int16); ok {
			// Convert numeric role to string
			roleMap := map[int16]string{
				0: "admin",
				1: "seller",
				2: "customer",
			}
			if roleStr, exists := roleMap[r]; exists {
				role = roleStr
			} else {
				role = "customer"
			}
		}

		// Try to get user from database
		user, err := m.repository.GetUserByID(userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// User doesn't exist, create new user
				newUser := &models.User{
					ID:        userID,
					Username:  username,
					Email:     email,
					Role:      role,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				err = m.repository.CreateUser(newUser)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to create user", nil, nil))
					return
				}

				user = newUser
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Database error", nil, nil))
				return
			}
		}

		// Set user in context
		c.Set("currentUser", *user)
		c.Set("userID", userID)
		c.Next()
	}
}
