package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/clients"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/common"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/transformers"

	// Import from common-protos
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/auth_service"
)

type AuthHandler struct {
	AuthClient *clients.AuthServiceClient
}

func NewAuthHandler(AuthClient *clients.AuthServiceClient) *AuthHandler {
	return &AuthHandler{AuthClient: AuthClient}
}

// Login
func (h *AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if reqBody.Email == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing email")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.AuthClient.Client.LoginUser(ctx, &pb.LoginRequest{Email: reqBody.Email, Password: reqBody.Password})
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.LoginRespJSON(grpcResp)

	token := httpResp["data"].(string)
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900, // 15 minutes in seconds to match exp claim
	})

	common.RespondWithJSON(w, http.StatusOK, httpResp)
}
