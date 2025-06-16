package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/clients"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/common"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/transformers"

	// Import from common-protos
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/auth_service"
)

type TokenClaims struct {
	UserID   string `json:"userID"`
	LastLog  string `json:"lastLog"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
	Type     string `json:"type"`
	HttpOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
}

type AuthHandler struct {
	AuthClient *clients.AuthServiceClient
}

func NewAuthHandler(AuthClient *clients.AuthServiceClient) *AuthHandler {
	return &AuthHandler{AuthClient: AuthClient}
}

// Validates and decodes JWT tokens from the auth service
func ValidateToken(tokenString string) (*TokenClaims, error) {
	// Empty check
	if tokenString == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Parse and validate
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Use the same secret key as auth service
		return []byte("Thunderbolts*"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Verify token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract and validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate required claims exist
	requiredClaims := []string{"userID", "lastLog", "iat", "exp", "type", "httpOnly", "secure"}
	for _, claim := range requiredClaims {
		if _, exists := claims[claim]; !exists {
			return nil, fmt.Errorf("missing required claim: %s", claim)
		}
	}

	// Validate token type
	if claims["type"] != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Validate token expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expiration claim")
	}
	if time.Now().Unix() > int64(exp) {
		return nil, fmt.Errorf("token has expired")
	}

	// TokenClaims to structured format
	tokenClaims := &TokenClaims{
		UserID:   claims["userID"].(string),
		LastLog:  claims["lastLog"].(string),
		Iat:      int64(claims["iat"].(float64)),
		Exp:      int64(claims["exp"].(float64)),
		Type:     claims["type"].(string),
		HttpOnly: claims["httpOnly"].(bool),
		Secure:   claims["secure"].(bool),
	}

	return tokenClaims, nil
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
	common.RespondWithJSON(w, http.StatusOK, httpResp)
}
