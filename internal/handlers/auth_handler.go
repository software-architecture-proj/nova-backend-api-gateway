package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/software-architecture-proj/nova-backend-auth-service/client"
	"github.com/software-architecture-proj/nova-backend-auth-service/transformer"
)

// Handler contains all the HTTP handlers and their dependencies
type Handler struct {
	authClient *client.AuthClient
}

// NewHandler creates a new Handler instance
func NewHandler(authClient *client.AuthClient) *Handler {
	return &Handler{
		authClient: authClient,
	}
}

// LoginRequest represents the HTTP request body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the HTTP response for login
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Email   string `json:"email"`
}

// LoginHandler handles the login HTTP request
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var httpReq transformer.HTTPLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&httpReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Transform HTTP request to gRPC request
	grpcReq := transformer.ToGRPCLoginRequest(&httpReq)

	// Call gRPC service
	grpcResp, err := h.authClient.LoginUser(r.Context(), grpcReq)
	if err != nil {
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// Transform gRPC response to HTTP response
	httpResp := transformer.ToHTTPLoginResponse(grpcResp)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(httpResp)
}
