package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/abhilash111/ecom/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Try session first
		sessionID, err := ctx.Cookie("session_id")
		if err == nil {
			valid, user, err := authService.ValidateSession(sessionID)
			fmt.Println("Session ID:", sessionID, "Valid:", valid, "User:", user, "Error:", err)
			if err == nil && valid && user != nil {
				fmt.Println("Session validated successfully", "User ID:", user.ID, "Role:", user.Role)
				ctx.Set("userID", user.ID)
				ctx.Set("userRole", string(user.Role))
				ctx.Next()
				return
			}
		}

		fmt.Println("Session validation failed, falling back to JWT")

		// Fall back to JWT
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bearer token required"})
			return
		}

		token, err := authService.ParseJWT(tokenString)
		fmt.Println("Token:", token, "Error:", err)
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		userID := uint(claims["id"].(float64))
		userRole := claims["role"].(string)
		ctx.Set("userID", userID)
		ctx.Set("userRole", userRole)

		ctx.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("userRole")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role information missing"})
			return
		}

		if userRole != requiredRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		ctx.Next()
	}
}
