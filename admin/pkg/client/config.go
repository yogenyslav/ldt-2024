package client

// GrpcClientConfig конфигурация клиента gRPC.
type GrpcClientConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
