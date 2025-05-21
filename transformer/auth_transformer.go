package transformer

import (
	"github.com/software-architecture-proj/nova-backend-auth-service/models"
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/auth_service"
)

// HTTPLoginRequest represents the HTTP request body for login
type HTTPLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HTTPLoginResponse represents the HTTP response for login
type HTTPLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Email   string `json:"email"`
}

// ToGRPCLoginRequest converts HTTP request to gRPC request
func ToGRPCLoginRequest(req *HTTPLoginRequest) *pb.LoginRequest {
	return &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

// ToHTTPLoginResponse converts gRPC response to HTTP response
func ToHTTPLoginResponse(resp *pb.LoginResponse) *HTTPLoginResponse {
	return &HTTPLoginResponse{
		Success: resp.Success,
		Message: resp.Message,
		Email:   resp.Email,
	}
}

// ToUserModel converts HTTP request to User model
func ToUserModel(req *HTTPLoginRequest) *models.User {
	return &models.User{
		Email:    req.Email,
		Password: req.Password,
	}
}
