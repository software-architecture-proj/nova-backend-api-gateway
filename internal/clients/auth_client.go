package clients

import (
	"log"

	// Import from common-protos
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	Client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthServiceClient(grpcHost string) (*AuthServiceClient, error) {
	conn, err := grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect to AuthService: %v", err)
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthServiceClient{Client: client, conn: conn}, nil
}

func (c *AuthServiceClient) CloseConnection() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
