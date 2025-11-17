package middleware

import (
	"app/domain/models"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"errors"
	"net/http"
	"os"
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
		token, err := jwt.ParseWithClaims(requestToken, &jwt_helpers.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid token signature", nil, nil))
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Token expired", nil, nil))
				return
			}

			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, err.Error(), nil, nil))
				return
			}
		}

		claims, ok := token.Claims.(*jwt_helpers.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid token claims", nil, nil))
			return
		}

		userID := helpers.ConvertStringToUint64(claims.UserID)
		username := claims.Username
		email := claims.Email

		// Handle role as either string or number
		var role string
		r := claims.Role
		// Convert numeric role to string
		roleMap := map[uint8]string{
			0: "admin",
			1: "seller",
			2: "customer",
		}
		if roleStr, exists := roleMap[r]; exists {
			role = roleStr
		} else {
			role = "customer"
		}

		// Try to get user from database
		u, err := m.repository.GetUserByID(userID)
		var user models.User
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

				if err = m.repository.CreateUser(newUser); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to create user", nil, nil))
					return
				}

				if newUser.Role == "admin" {
					if err = m.repository.CreateAdminAvailability(&models.AdminAvailability{
						AdminID: uint8(newUser.ID),
					}); err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to create admin availability", nil, nil))
						return
					}
				}

				user = *newUser
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Database error", nil, nil))
				return
			}
		} else {
			// existing user
			if u == nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Database error", nil, nil))
				return
			}
			user = *u
		}

		// set model user (not jwt claims) into context
		c.Set("userData", user)
		c.Next()
	}
}
