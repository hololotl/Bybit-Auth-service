package auth

import (
	"Bybit_Pet_Project/authServer/grpcFile"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerAPI struct {
	grpcFile.UnimplementedAuthServiceServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error)
}

func RegisterServer(server *grpc.Server, auth Auth) {
	grpcFile.RegisterAuthServiceServer(server, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Login(ctx context.Context, in *grpcFile.LoginRequest) (*grpcFile.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "AppId is required")
	}
	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {

	}
	return &grpcFile.LoginResponse{Token: token}, nil
}
func (s *ServerAPI) Register(ctx context.Context, in *grpcFile.RegisterRequest) (*grpcFile.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}
	uid, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {

	}
	return &grpcFile.RegisterResponse{UserId: uid}, nil
}
