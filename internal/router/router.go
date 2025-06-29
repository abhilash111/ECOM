package router

import (
	"database/sql"

	"github.com/abhilash111/ecom/internal/products"
	user "github.com/abhilash111/ecom/internal/users"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	// user routes
	userStore := user.NewStore(db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(api)

	// product routes
	productStore := products.NewStore(db)
	productHandler := products.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(api)

	return r
}
