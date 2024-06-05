package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClient is a wrapper around grpc.ClientConn
type GrpcClient struct {
	conn *grpc.ClientConn
}

// NewGrpcClient creates a new GrpcClient
func NewGrpcClient(cfg *GrpcClientConfig) (*GrpcClient, error) {
	var grpcOpts []grpc.DialOption
	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.NewClient(addr, grpcOpts...)
	if err != nil {
		return nil, err
	}

	return &GrpcClient{conn: conn}, nil
}

// Close closes the grpc connection
func (c *GrpcClient) Close() error {
	return c.conn.Close()
}

// GetConn returns the grpc connection
func (c *GrpcClient) GetConn() *grpc.ClientConn {
	return c.conn
}
