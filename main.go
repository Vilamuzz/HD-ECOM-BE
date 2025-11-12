package main

import (
	"app/app/handlers"
	"app/app/middleware"
	"app/app/repositories"
	"app/app/services"
	"app/docs"
	"app/domain"
	"app/helpers"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	helpers.LoadEnv()
}

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	docs.SwaggerInfo.Title = "Helpdesk E-Commerce API"
	docs.SwaggerInfo.Description = "Api documentation for Helpdesk E-Commerce Application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + os.Getenv("APP_PORT")
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	db := helpers.ConnectDB()
	helpers.MigrateDB(db, domain.GetAllModels()...)
	repo := repositories.NewAppRepository(db)
	hub := services.NewHub()
	service := services.NewAppService(repo, hub)
	go service.Run()
	middleware := middleware.NewAppMiddleware(repo)
	ginEngine := gin.Default()

	// Add CORS middleware
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))

	handlers.App(service, repo, ginEngine, middleware)

	ginEngine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"message": "Hello World!",
		})
	})

	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("APP_PORT")
	ginEngine.Run(":" + port)
}
