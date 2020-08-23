package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	envoy_api_v3_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/reflection"
)

type EnvoyV3Server struct {
}

func NewEnvoyV3Server() EnvoyV3Server {
	return EnvoyV3Server{}
}

func (s EnvoyV3Server) Run(port int) {
	address := fmt.Sprintf(":%d", port)

	log.Info().Msgf("Listening for gRPC on %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to listen on address: %s", address)
	}

	log.Info().Msgf("Starting gRPC server on %s...", address)
	grpcServer := newServer()

	envoy_service_auth_v3.RegisterAuthorizationServer(grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Info().Msgf("Registered Envoy authz gRPC services on %s...", address)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msgf("Failed to serve gRPC on address: %s", address)
	}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s EnvoyV3Server) Check(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	tr := global.Tracer("")
	_, span := tr.Start(ctx, "Check")
	defer span.End()

	log.Info().Msg("ENVOY GATEWAY HIT!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	m := jsonpb.Marshaler{}
	pbs, _ := m.MarshalToString(req)
	log.Info().Msg(pbs)

	authorization := req.Attributes.Request.Http.Headers["authorization"]

	log.Info().Msgf("Authorization header is: %s", authorization)

	extracted := strings.Fields(authorization)

	if len(extracted) != 2 || extracted[0] != "Bearer" {
		log.Err(ErrMalformedAuthHeader).Msg("Malformed auth header")
		return envoyV3ErrResponse(ErrMalformedAuthHeader), nil
	}

	token := extracted[1]
	if token != "kevin" {
		err := errors.New("bad token")
		log.Err(err).Msg("Bad token")
		return envoyV3ErrResponse(err), nil
	}

	log.Info().Msg("Good token")
	return &envoy_service_auth_v3.CheckResponse{
		HttpResponse: &envoy_service_auth_v3.CheckResponse_OkResponse{
			OkResponse: &envoy_service_auth_v3.OkHttpResponse{
				Headers: []*envoy_api_v3_core.HeaderValueOption{
					{
						Append: &wrappers.BoolValue{Value: false},
						Header: &envoy_api_v3_core.HeaderValue{
							// For a successful request, the authorization server sets the
							// x-current-user value.
							Key:   "x-current-user",
							Value: "kevin",
						},
					},
				},
			},
		},
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
	}, nil
}

func envoyV3ErrResponse(err error) *envoy_service_auth_v3.CheckResponse {
	s := &status.Status{
		Code: int32(code.Code_PERMISSION_DENIED),
	}
	out := &envoy_service_auth_v3.CheckResponse{
		Status: s,
		HttpResponse: &envoy_service_auth_v3.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
				Status: &envoy_type_v3.HttpStatus{
					Code: envoy_type_v3.StatusCode_Forbidden,
				},
				Headers: nil,
				Body:    "",
			},
		},
	}

	// Set the error message
	if err != nil {
		out.GetDeniedResponse().Body = err.Error()
	}

	return out
}
