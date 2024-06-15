package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClient интерфейс для работы с gRPC клиентом.
type GrpcClient interface {
	GetConn() *grpc.ClientConn
	Close() error
}

// GrpcClient обертка над gRPC клиентом.
type grpcClient struct {
	conn *grpc.ClientConn
}

// NewGrpcClient создает новый gRPC клиент.
func NewGrpcClient(cfg *GrpcClientConfig) (GrpcClient, error) {
	var grpcOpts []grpc.DialOption
	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.NewClient(addr, grpcOpts...)
	if err != nil {
		return nil, err
	}

	return &grpcClient{conn: conn}, nil
}

// Close закрывает gRPC клиент.
func (c *grpcClient) Close() error {
	return c.conn.Close()
}

// GetConn получить gRPC соединение.
func (c *grpcClient) GetConn() *grpc.ClientConn {
	return c.conn
}
