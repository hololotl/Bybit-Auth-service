package main

import (
	"Bybit_Pet_Project/internal/app"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	time2 "time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	log := setupLogger("1")
	time := time2.Duration(100000000000000)
	conString := "user=postgres password=2005 dbname=petProjBybit sslmode=disable"
	application := app.New(log, 8000, conString, time)
	application.GRPCServer.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
