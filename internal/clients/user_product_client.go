package clients

import (
	"log"

	// Import from common-protos
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/user_product_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserProductServiceClient struct {
	Client pb.UserProductServiceClient
	conn   *grpc.ClientConn
}

func NewUserProductServiceClient(grpcHost string) (*UserProductServiceClient, error) {
	conn, err := grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect to UserProductService: %v", err)
		return nil, err
	}

	client := pb.NewUserProductServiceClient(conn)
	return &UserProductServiceClient{Client: client, conn: conn}, nil
}

func (c *UserProductServiceClient) CloseConnection() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
