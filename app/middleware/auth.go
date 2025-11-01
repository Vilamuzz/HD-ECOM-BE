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
			log.Printf("[AUTH] Token from query parameter: %s", requestToken)
		} else {
			log.Printf("[AUTH] Token from Authorization header: %s", requestToken)
		}

		if requestToken == "" {
			log.Println("[AUTH] Missing token - no token in header or query")
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Missing token", nil, nil))
			return
		}

		// Remove "Bearer " prefix if present
		if strings.HasPrefix(strings.ToLower(requestToken), "bearer ") {
			requestToken = strings.TrimSpace(requestToken[7:])
			log.Printf("[AUTH] Token after removing 'Bearer ': %s", requestToken)
		}

		// Parse and validate JWT
		token, err := jwt.Parse(requestToken, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			log.Printf("[AUTH] JWT parse error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid token", nil, nil))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("[AUTH] Failed to parse token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid token claims", nil, nil))
			return
		}

		log.Printf("[AUTH] Token claims: %+v", claims)

		// Check expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			log.Println("[AUTH] Token missing exp claim")
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Token missing exp", nil, nil))
			return
		}
		if float64(time.Now().Unix()) > exp {
			log.Printf("[AUTH] Token expired. Now: %d, Exp: %d", time.Now().Unix(), int64(exp))
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Token expired", nil, nil))
			return
		}

		// Extract user_id from claims
		var userIDStr string

		// Try different claim formats
		if uid, ok := claims["user_id"].(string); ok {
			userIDStr = uid
		} else if uid, ok := claims["sub"].(string); ok {
			userIDStr = uid
		} else if uid, ok := claims["user_id"].(float64); ok {
			userIDStr = strconv.FormatInt(int64(uid), 10)
		} else if uid, ok := claims["sub"].(float64); ok {
			userIDStr = strconv.FormatInt(int64(uid), 10)
		} else {
			log.Printf("[AUTH] No valid user_id found in claims. Available claims: %+v", claims)
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid user_id in token", nil, nil))
			return
		}

		log.Printf("[AUTH] Extracted userID: %s", userIDStr)

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Printf("[AUTH] Failed to parse userID '%s': %v", userIDStr, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse(http.StatusUnauthorized, "Invalid user_id format", nil, nil))
			return
		}

		// Extract other user info
		username, _ := claims["username"].(string)
		email, _ := claims["email"].(string)

		// Handle role as either string or number
		var role string
		if r, ok := claims["role"].(string); ok {
			role = r
		} else if r, ok := claims["role"].(float64); ok {
			// Convert numeric role to string
			roleMap := map[float64]string{
				0: "admin",
				1: "seller",
				2: "buyer",
			}
			if roleStr, exists := roleMap[r]; exists {
				role = roleStr
			} else {
				role = "customer" // default
			}
		}

		log.Printf("[AUTH] User info - ID: %d, Username: %s, Email: %s, Role: %s", userID, username, email, role)

		// Try to get user from database
		user, err := m.repository.GetUserByID(userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("[AUTH] User %d not found in DB, creating new user", userID)

				// User doesn't exist, create new user
				newUser := &models.User{
					ID:        userID,
					Username:  username,
					Email:     email,
					Role:      role,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				// Set default role if not provided
				if newUser.Role == "" {
					newUser.Role = "customer"
				}

				// Set default username if not provided
				if newUser.Username == "" {
					newUser.Username = email
				}

				err = m.repository.CreateUser(newUser)
				if err != nil {
					log.Printf("[AUTH] Error creating user: %v", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to create user", nil, nil))
					return
				}

				user = newUser
				log.Printf("[AUTH] Created new user: ID=%d, Username=%s, Email=%s, Role=%s", user.ID, user.Username, user.Email, user.Role)
			} else {
				log.Printf("[AUTH] Database error getting user: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Database error", nil, nil))
				return
			}
		} else {
			log.Printf("[AUTH] Found existing user: ID=%d, Username=%s", user.ID, user.Username)
		}

		// Set user in context
		c.Set("currentUser", *user)
		c.Set("userID", userID)

		log.Printf("[AUTH] Successfully authenticated user %d", userID)
		c.Next()
	}
}
