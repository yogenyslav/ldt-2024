package client

// GrpcClientConfig is the configuration for the grpc client
type GrpcClientConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
