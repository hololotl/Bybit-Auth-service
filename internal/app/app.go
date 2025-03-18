package app

import (
	grpcapp "Bybit_Pet_Project/internal/app/grpc"
	"Bybit_Pet_Project/internal/services/auth"
	postgr "Bybit_Pet_Project/internal/storage/sql"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagepath string, tokenTTL time.Duration) *App {
	storage, err := postgr.New(storagepath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{GRPCServer: grpcApp}
}
