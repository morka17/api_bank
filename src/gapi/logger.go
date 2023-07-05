package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	result, err := handler(ctx, req)
	startTime := time.Now()
	duration := time.Since(startTime)
	
	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	
	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("Protocol", "grpc").
		Str("method", info.FullMethod). 
		Int("status_code", int(statusCode)). 
		Str("status_text", statusCode.String()).
		Dur("duration", duration). 
		Msg("received a gRPC request")
	
	
	return result, err 
	
}