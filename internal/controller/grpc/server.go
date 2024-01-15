package grpc

import (
	"jastip-app/internal/usecase"

	"google.golang.org/grpc"
)

func NewGrpcServer(s *grpc.Server, uc *usecase.UseCase) *grpc.Server {
	return s
}
