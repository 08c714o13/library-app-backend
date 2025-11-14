package middlewares

import (
	"os"
	"time"

	"library/internal/app/domain"
	"library/internal/app/services"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v5"
)

const (
	identityKey = "email"
	roleKey     = "role"
)

func InitAuthJWT(adminService services.AdminService, librarianService services.LibrarianService) (*jwt.GinJWTMiddleware, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key"
	}

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:           "library zone",
		Key:             []byte(jwtSecret),
		Timeout:         time.Hour * 24,
		MaxRefresh:      time.Hour * 24,
		IdentityKey:     identityKey,
		PayloadFunc:     payloadFunc(),
		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(adminService, librarianService),
		Authorizer:      authorizer(),
		Unauthorized:    unauthorized(),
		TimeFunc:        time.Now,
		SendCookie:      true,
	}

	if err := authMiddleware.MiddlewareInit(); err != nil {
		return nil, err
	}

	return authMiddleware, nil
}

func payloadFunc() func(data any) jwt2.MapClaims {
	return func(data any) jwt2.MapClaims {
		if v, ok := data.(*domain.User); ok {
			return jwt2.MapClaims{
				identityKey: v.Email,
				roleKey:     v.Role,
			}
		}
		return jwt2.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) any {
	return func(c *gin.Context) any {
		claims := jwt.ExtractClaims(c)
		return &domain.User{
			Email: claims[identityKey].(string),
			Role:  claims[roleKey].(string),
		}
	}
}

func authenticator(adminService services.AdminService, librarianService services.LibrarianService) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var loginVals domain.LoginPayload
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		email := loginVals.Email
		password := loginVals.Password

		ok, err := adminService.Login(email, password)
		if err == nil && ok {
			return &domain.User{
				Email: email,
				Role:  "admin",
			}, nil
		}

		ok, err = librarianService.Login(email, password)
		if err == nil && ok {
			return &domain.User{
				Email: email,
				Role:  "librarian",
			}, nil
		}

		return nil, jwt.ErrFailedAuthentication
	}
}

func authorizer() func(c *gin.Context, data any) bool {
	return func(c *gin.Context, data any) bool {
		if v, ok := data.(*domain.User); ok {
			return v.Role == "admin" || v.Role == "librarian"
		}
		return false
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		role, exists := claims[roleKey].(string)
		if !exists {
			c.JSON(403, gin.H{"error": "no role found in token"})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

// func handleNoRoute() func(c *gin.Context) {
//	return func(c *gin.Context) {
// 		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
// 	}
// }
