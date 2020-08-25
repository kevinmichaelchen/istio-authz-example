package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_service_auth_v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoy_type_v2 "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/reflection"
)

type EnvoyV2Server struct {
}

func (s EnvoyV2Server) String() string {
	return "EnvoyV2Server"
}

func NewEnvoyV2Server() EnvoyV2Server {
	return EnvoyV2Server{}
}

func (s EnvoyV2Server) Run(port int) {
	address := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to listen on address: %s", address)
	}

	log.Info().Msgf("Starting %s on %s...", s, address)
	grpcServer := newServer()

	envoy_service_auth_v2.RegisterAuthorizationServer(grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Info().Msgf("Registered Envoy authz gRPC services on %s...", address)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msgf("Failed to serve gRPC on address: %s", address)
	}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s EnvoyV2Server) Check(ctx context.Context, req *envoy_service_auth_v2.CheckRequest) (*envoy_service_auth_v2.CheckResponse, error) {
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
		return envoyV2ErrResponse(ErrMalformedAuthHeader), nil
	}

	token := extracted[1]
	if token != "kevin" {
		err := errors.New("bad token")
		log.Err(err).Msg("Bad token")
		return envoyV2ErrResponse(err), nil
	}

	log.Info().Msg("Good token")
	return &envoy_service_auth_v2.CheckResponse{
		HttpResponse: &envoy_service_auth_v2.CheckResponse_OkResponse{
			OkResponse: &envoy_service_auth_v2.OkHttpResponse{
				Headers: []*envoy_api_v2_core.HeaderValueOption{
					{
						Append: &wrappers.BoolValue{Value: false},
						Header: &envoy_api_v2_core.HeaderValue{
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

func envoyV2ErrResponse(err error) *envoy_service_auth_v2.CheckResponse {
	s := &status.Status{
		Code: int32(code.Code_PERMISSION_DENIED),
	}
	out := &envoy_service_auth_v2.CheckResponse{
		Status: s,
		HttpResponse: &envoy_service_auth_v2.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v2.DeniedHttpResponse{
				Status: &envoy_type_v2.HttpStatus{
					Code: envoy_type_v2.StatusCode_Forbidden,
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
