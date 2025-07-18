package clients

import (
	"log"

	// Import from common-protos
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/transaction_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TransactionServiceClient struct {
	Client pb.TransactionServiceClient
	conn   *grpc.ClientConn
}

func NewTransactionServiceClient(grpcHost string) (*TransactionServiceClient, error) {
	conn, err := grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("did not connect to TransactionService: ", err)
		return nil, err
	}

	client := pb.NewTransactionServiceClient(conn)
	return &TransactionServiceClient{Client: client, conn: conn}, nil
}

func (c *TransactionServiceClient) CloseConnection() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
