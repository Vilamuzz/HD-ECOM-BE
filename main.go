package main

import (
	"app/app/handlers"
	"app/app/middleware"
	"app/app/repositories"
	"app/app/repositories/s3"
	"app/app/services"
	"app/docs"
	"app/domain"
	"app/helpers"
	"log"
	"net/http"
	"os"
	"strconv"
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
	docs.SwaggerInfo.Host = os.Getenv("BASE_API_URL")
	docs.SwaggerInfo.BasePath = "/api"
	scheme := os.Getenv("SWAGGER_SCHEME")
	if scheme == "" {
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		docs.SwaggerInfo.Schemes = []string{scheme}
	}

	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "5"
	}
	timeout, _ := strconv.Atoi(timeoutStr)
	timeoutContext := time.Duration(timeout) * time.Second

	db := helpers.ConnectDB()
	helpers.MigrateDB(db, domain.GetAllModels()...)
	repo := repositories.NewAppRepository(db)

	// Add S3 repository initialization
	s3Repo := s3.NewS3Repository(timeoutContext)

	hub := services.NewHub()
	service := services.NewAppService(services.DBInjection{
		Repo:   repo,
		S3Repo: s3Repo,
	}, hub, timeoutContext)
	go service.Run()
	middleware := middleware.NewAppMiddleware(repo)
	ginEngine := gin.Default()

	// Add CORS middleware
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", os.Getenv("APP_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))

	handlers.App(service, ginEngine, middleware)

	ginEngine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"message": "Hello World!",
		})
	})

	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start cleanup job for expired messages
	go startMessageCleanupJob(repo)

	port := os.Getenv("APP_PORT")
	ginEngine.Run(":" + port)
}

// startMessageCleanupJob runs daily to permanently delete messages past their purge date
func startMessageCleanupJob(repo domain.AppRepository) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := repo.PermanentlyDeleteExpiredMessages(); err != nil {
			log.Printf("Error cleaning up expired messages: %v", err)
		} else {
			log.Println("Successfully cleaned up expired messages")
		}
	}
}
