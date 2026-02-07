package main

import (
	"context"
	v1 "emias_printer/internal/transport/http/v1"
	"emias_printer/pkg/config"
	"emias_printer/pkg/logger"
	"emias_printer/pkg/printer"
	"fmt"
	"log"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	ctx, err := logger.NewLogger(ctx)
	if err != nil {
		log.Fatalf("cannot init logger: %v", err)
	}

	cfg, err := config.NewConfig("config/.env")
	if err != nil {
		logger.GetLoggerFromContext(ctx).Fatal(ctx, "cannot parse config", zap.Error(err))
	}
	println(cfg.Delay)
	mp := printer.NewPrinterManipulator(cfg.Delay)

	srv := v1.NewServer(cfg.Port)  
	err = srv.RegisterHandlers(ctx, mp)
	if err != nil {
		log.Fatalf("cannot register handlers: %v", err)
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, fmt.Sprintf("starting web server on port: %v", cfg.Port))
	if err = srv.Start(); err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
	 
}


// @title My API
// @version 1.0
// @description Example API
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

