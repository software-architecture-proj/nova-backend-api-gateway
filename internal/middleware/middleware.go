package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	// Import from common-protos
)

type TokenClaims struct {
	UserID   string `json:"userID"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	LastLog  string `json:"lastLog"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
	Type     string `json:"type"`
	HttpOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
}

type MiddlewareInterface interface {
	AuthToken(handler http.Handler) http.Handler
}

func NewMiddleware() MiddlewareInterface {
	return &TokenClaims{}
}

func (m *TokenClaims) AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from cookie first
		tokenString, err := r.Cookie("accessToken")
		var tokenValue string
		if err == nil && tokenString.Value != "" {
			tokenValue = tokenString.Value
		} else {
			// If no cookie, try Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization is required", http.StatusUnauthorized)
				return
			}
			// Check if it's a Bearer token
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenValue = authHeader[7:]
			} else {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}
		}

		// Validate and decode the token
		tokenClaims, err := validateToken(tokenValue)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Store claims in context for downstream handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "tokenClaims", tokenClaims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r) // Call the next handler
	})
}

// Validates and decodes JWT tokens from the auth service
func validateToken(tokenString string) (*TokenClaims, error) {
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
		Email:    claims["email"].(string),
		Username: claims["username"].(string),
		Phone:    claims["phone"].(string),
		LastLog:  claims["lastLog"].(string),
		Iat:      int64(claims["iat"].(float64)),
		Exp:      int64(claims["exp"].(float64)),
		Type:     claims["type"].(string),
		HttpOnly: claims["httpOnly"].(bool),
		Secure:   claims["secure"].(bool),
	}

	return tokenClaims, nil
}
