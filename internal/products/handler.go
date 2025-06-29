package products

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/abhilash111/ecom/internal/auth"
	"github.com/abhilash111/ecom/internal/types"
	"github.com/abhilash111/ecom/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/products", h.handleGetProducts)
	router.GET("/products/:productID", h.handleGetProduct)

	authGroup := router.Group("/products")
	authGroup.Use(auth.JWTAuthMiddleware(h.userStore)) // Apply JWT auth middleware to this group
	// admin routes
	authGroup.POST("", h.handleCreateProduct)
	// router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	// router.HandleFunc("/products/{productID}", h.handleGetProduct).Methods(http.MethodGet)

	// // admin routes
	// router.HandleFunc("/products", auth.WithJWTAuth(h.handleCreateProduct, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleGetProducts(c *gin.Context) {
	products, err := h.store.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) handleGetProduct(c *gin.Context) {
	productIDStr := c.Param("productID")
	productID, err := strconv.Atoi(productIDStr)

	if err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	product, err := h.store.GetProductByID(productID)
	fmt.Printf("Product ID: %d, Product: %+v\n", productID, product)
	fmt.Printf("Error: %v\n", err)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(c.Writer, http.StatusOK, product)
}

func (h *Handler) handleCreateProduct(c *gin.Context) {
	var product types.CreateProductPayload
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.CreateProduct(product)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(c.Writer, http.StatusCreated, product)
}
