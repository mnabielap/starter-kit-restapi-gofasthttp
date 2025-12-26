package middleware

import (
	"strings"

	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/service"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/pkg/utils"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// AuthMiddleware handles JWT verification and Role-Based Access Control (RBAC)
func AuthMiddleware(tokenService *service.TokenService, requiredRoles ...string) routing.Handler {
	return func(c *routing.Context) error {
		// 1. Check Authorization Header
		authHeader := string(c.Request.Header.Peek("Authorization"))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteError(c.RequestCtx, fasthttp.StatusUnauthorized, "Missing or invalid Authorization header")
			c.Abort()
			return nil
		}

		// 2. Verify Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := tokenService.VerifyToken(tokenString)

		// Validate token validity and type (must be access token)
		if err != nil || claims.Type != model.TokenTypeAccess {
			utils.WriteError(c.RequestCtx, fasthttp.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return nil
		}

		// 3. RBAC Check (If roles are specified)
		if len(requiredRoles) > 0 {
			hasRole := false
			for _, role := range requiredRoles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}
			if !hasRole {
				utils.WriteError(c.RequestCtx, fasthttp.StatusForbidden, "Forbidden: Insufficient permissions")
				c.Abort()
				return nil
			}
		}

		// 4. Store user info in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		return c.Next()
	}
}