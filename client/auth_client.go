package client

import (
	"context"

	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/auth_service"
	"google.golang.org/grpc"
)

// AuthClient handles gRPC communication with the auth service
type AuthClient struct {
	client pb.AuthServiceClient
}

// NewAuthClient creates a new AuthClient instance
func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		client: pb.NewAuthServiceClient(conn),
	}
}

// LoginUser sends a login request to the auth service
func (c *AuthClient) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return c.client.LoginUser(ctx, req)
}
