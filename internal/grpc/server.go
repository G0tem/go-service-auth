package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/G0tem/go-servise-auth/internal/config"
	"github.com/G0tem/go-servise-auth/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// AuthServer реализует gRPC сервер для авторизации
type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	cfg *config.Config
}

// NewAuthServer создает новый экземпляр gRPC сервера
func NewAuthServer(cfg *config.Config) *AuthServer {
	return &AuthServer{
		cfg: cfg,
	}
}

// GetTestData возвращает тестовые данные для проверки gRPC запроса
func (s *AuthServer) GetTestData(ctx context.Context, req *proto.GetTestDataRequest) (*proto.GetTestDataResponse, error) {
	log.Info().
		Str("message", req.Message).
		Msg("Received GetTestData gRPC request")

	return &proto.GetTestDataResponse{
		Message:   fmt.Sprintf("Hello from auth service! Your message: %s", req.Message),
		Status:    200,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// GetUserInfo возвращает информацию о пользователе по ID
func (s *AuthServer) GetUserInfo(ctx context.Context, req *proto.GetUserInfoRequest) (*proto.GetUserInfoResponse, error) {
	log.Info().
		Str("user_id", req.UserId).
		Msg("Received GetUserInfo gRPC request")

	// Пример тестовых данных
	return &proto.GetUserInfoResponse{
		UserId:   req.UserId,
		Email:    "test@example.com",
		Username: "testuser",
		IsActive: true,
	}, nil
}

// StartGrpcServer запускает gRPC сервер
func StartGrpcServer(cfg *config.Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authServer := NewAuthServer(cfg)
	proto.RegisterAuthServiceServer(s, authServer)

	log.Info().Msgf("gRPC server listening on port %d", cfg.GrpcPort)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
