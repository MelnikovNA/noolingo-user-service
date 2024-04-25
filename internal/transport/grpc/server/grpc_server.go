package grpcserver

import (
	"fmt"
	"net"

	"github.com/noolingo/proto/codegen/go/noolingo"
	"github.com/noolingo/user-service/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	host    string
	port    string
	server  *grpc.Server
	service *service.Services
	logger  *logrus.Logger
}

func New(host string, port string, service *service.Services, logger *logrus.Logger) *Server {
	return &Server{
		host:    host,
		port:    port,
		server:  grpc.NewServer(),
		service: service,
		logger:  logger,
	}
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return fmt.Errorf("can't start listening addr: %w", err)
	}
	grpc_health_v1.RegisterHealthServer(s.server, health.NewServer())
	noolingo.RegisterUserServer(s.server, newUserServer(s.logger, s.service))
	err = s.server.Serve(lis)

	if err != nil {
		return fmt.Errorf("can't start listening addr for grpc server: %w", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
