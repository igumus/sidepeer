package grpc

import (
	"errors"
	"fmt"
)

const defaultGrpcHost = "0.0.0.0"

const defaultGrpcPort = 4000

const _minPort = 0

const _maxPort = 65535

// Captures/Represents GRPC server configuration information.
type grpcServerConfig struct {
	host string
	port int
}

func (s *grpcServerConfig) String() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

// Validate - validates grpcServerConfig instance.
func validate(s *grpcServerConfig) error {
	if len(s.host) == 0 {
		return errors.New("[sidepeer] grpc server configuration failed: host not specified")
	}
	if s.port < _minPort || s.port >= _maxPort {
		return errors.New("[sidepeer] grpc server configuration failed: port is not valid")
	}
	return nil
}

func defaultGrpcServerConfig() *grpcServerConfig {
	return &grpcServerConfig{
		host: defaultGrpcHost,
		port: defaultGrpcPort,
	}
}

func applyOptions(opts ...GrpcServerOption) (*grpcServerConfig, error) {
	ret := defaultGrpcServerConfig()
	for _, opt := range opts {
		opt(ret)
	}
	return ret, validate(ret)
}

type GrpcServerOption func(*grpcServerConfig)

func WithPort(p int) GrpcServerOption {
	return func(bsc *grpcServerConfig) {
		bsc.port = p
	}
}
