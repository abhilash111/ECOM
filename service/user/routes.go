package user

import (
	"net/http"

	"github.com/abhilash111/ecom/types"
	"github.com/abhilash111/ecom/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewUserHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Get JSON Payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
}
