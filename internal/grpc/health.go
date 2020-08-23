package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/reflection"

	"github.com/rs/zerolog/log"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct{}

func NewHealthServer() HealthServer {
	return HealthServer{}
}
func (s HealthServer) Check(ctx context.Context, request *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

// rpc Watch(HealthCheckRequest) returns (stream HealthCheckResponse);
func (s HealthServer) Watch(request *health.HealthCheckRequest, stream health.Health_WatchServer) error {
	// Send healthy
	if err := stream.Send(&health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}); err != nil {
		log.Error().Err(err).Msg("failed to send response into gRPC stream")
	}

	return nil
}

func (s HealthServer) Run(port int) {
	address := fmt.Sprintf(":%d", port)

	log.Info().Msgf("Listening for gRPC on %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to listen on address: %s", address)
	}

	log.Info().Msgf("Starting gRPC server on %s...", address)
	grpcServer := newServer()

	health.RegisterHealthServer(grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Info().Msgf("Registered Envoy authz gRPC services on %s...", address)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msgf("Failed to serve gRPC on address: %s", address)
	}
}
