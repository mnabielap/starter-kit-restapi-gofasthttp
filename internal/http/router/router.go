package router

import (
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/http/handler"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/http/middleware"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/service"
	routing "github.com/qiangxue/fasthttp-routing"
	fasthttpSwagger "github.com/swaggo/fasthttp-swagger"
	"github.com/valyala/fasthttp"
)

// adaptHandler converts a standard fasthttp.RequestHandler to routing.Handler
func adaptHandler(h fasthttp.RequestHandler) routing.Handler {
	return func(c *routing.Context) error {
		h(c.RequestCtx)
		return nil
	}
}

func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tokenService *service.TokenService,
) *routing.Router {
	router := routing.New()

	router.Use(middleware.Logger)

	// Swagger Docs
	router.Get("/swagger/*", adaptHandler(fasthttpSwagger.WrapHandler(fasthttpSwagger.DeepLinking(true))))

	v1 := router.Group("/v1")

	// --- Auth Routes (Public) ---
	auth := v1.Group("/auth")
	auth.Post("/register", adaptHandler(authHandler.Register))
	auth.Post("/login", adaptHandler(authHandler.Login))
	auth.Post("/logout", adaptHandler(authHandler.Logout))
	auth.Post("/refresh-tokens", adaptHandler(authHandler.RefreshTokens))

	// --- User Routes (Protected: Admin Only) ---
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(tokenService, "admin"))
	
	users.Post("", adaptHandler(userHandler.CreateUser))
	users.Get("", adaptHandler(userHandler.GetUsers))
	users.Get("/<userId>", userHandler.GetUser)
	users.Patch("/<userId>", userHandler.UpdateUser)
	users.Delete("/<userId>", userHandler.DeleteUser)

	return router
}