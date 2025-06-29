package auth

import (
	"log"
	"net/http"
	"strconv"

	"github.com/abhilash111/ecom/internal/types"
	"github.com/abhilash111/ecom/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(store types.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := utils.GetTokenFromRequest(c.Request)

		token, err := validateJWT(tokenStr)
		if err != nil || !token.Valid {
			log.Println("invalid or missing token:", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		strID := claims["userID"].(string)
		userID, err := strconv.Atoi(strID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid user ID"})
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user not found"})
			return
		}

		c.Set("userID", u.ID)
		c.Next()
	}
}
