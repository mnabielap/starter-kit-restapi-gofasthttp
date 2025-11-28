package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/config"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/database"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/http/handler"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/http/router"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/repository"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/service"
	"github.com/valyala/fasthttp"

	// Import generated docs so swagger can load them
	_ "github.com/mnabielap/starter-kit-restapi-gofasthttp/docs"
)

// @title           Go FastHTTP Starter Kit API
// @version         1.0
// @description     This is a sample RESTful API using Go, FastHTTP, GORM, and Clean Architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT

// @host            localhost:3000
// @BasePath        /v1
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load Configuration
	config.LoadConfig()
	log.Println("✅ Configuration loaded.")

	// 2. Connect to Database
	db := database.ConnectDB(config.AppConfig)

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// 4. Initialize Services
	tokenService := service.NewTokenService(tokenRepo, config.AppConfig)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userService, tokenService)

	// 5. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// 6. Setup Router
	appRouter := router.SetupRouter(authHandler, userHandler, tokenService)
	log.Println("✅ API router initialized.")

	// 7. Start Server
	serverAddr := fmt.Sprintf(":%d", config.AppConfig.Port)
	server := &fasthttp.Server{
		Handler: appRouter.HandleRequest,
		Name:    "GoFastHTTP-Starter-Kit",
	}

	go func() {
		log.Printf("🚀 Server starting on http://localhost%s", serverAddr)
		log.Printf("📄 Swagger docs available at http://localhost%s/swagger/index.html", serverAddr)
		if err := server.ListenAndServe(serverAddr); err != nil {
			log.Fatalf("❌ Server error: %s", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("🚦 Shutting down server...")
	if err := server.Shutdown(); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}
	log.Println("👋 Server stopped.")
}