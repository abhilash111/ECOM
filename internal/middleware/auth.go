package middleware

import (
	"net/http"
	"strings"

	"github.com/abhilash111/ecom/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Set user information in context
		ctx.Set("userID", claims["id"])
		ctx.Set("userEmail", claims["email"])
		ctx.Set("userRole", claims["role"])

		ctx.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("userRole")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid role type"})
			return
		}

		hasAccess := false
		for _, role := range roles {
			if roleStr == role {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		ctx.Next()
	}
}
