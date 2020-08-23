package grpc

import (
	"strings"

	"github.com/kevinmichaelchen/istio-authz-example/internal/configuration"

	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tags"
	grpcGo "google.golang.org/grpc"
)

func newServer() *grpcGo.Server {
	// Logger is used, allowing pre-definition of certain fields by the user.
	logger := configuration.GetLogger()
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []logging.Option{
		logging.WithDecider(func(fullMethodName string, err error) bool {
			// Don't log health gRPC endpoint
			if strings.HasSuffix(fullMethodName, "/Check") {
				return false
			}
			return true
		}),
	}

	return grpcGo.NewServer(
		middleware.WithUnaryServerChain(
			tags.UnaryServerInterceptor(),
			// TODO wait until they use OpenTelemetry
			//opentracing.UnaryServerInterceptor(),
			//prometheus.UnaryServerInterceptor,
			logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger), opts...),
			recovery.UnaryServerInterceptor(),
		),
		middleware.WithStreamServerChain(
			tags.StreamServerInterceptor(),
			// TODO wait until they use OpenTelemetry
			//opentracing.StreamServerInterceptor(),
			//prometheus.StreamServerInterceptor,
			logging.StreamServerInterceptor(grpczerolog.InterceptorLogger(logger), opts...),
			recovery.StreamServerInterceptor(),
		),
	)
}
