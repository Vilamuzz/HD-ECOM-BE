package middleware

import (
    "app/domain/models"
    "github.com/gin-gonic/gin"
)

type AppMiddleware interface {
    Auth() gin.HandlerFunc
    RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc
    RequireAdminOrSupport() gin.HandlerFunc
}