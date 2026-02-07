package v1

import (
	"context"
	"emias_printer/internal/transport/http/v1/handlers"
	"emias_printer/pkg/logger"
	"fmt"
	"net/http"
	_"emias_printer/docs"
	"time"
	"emias_printer/pkg/printer"

	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	defaultHeaderTimeout = time.Second * 5
)

type Server struct {
	srv *http.Server
}

func NewServer(port int) *Server {
	srv := http.Server{
		Addr:        fmt.Sprintf(":%v", port),
		Handler:     nil,
		ReadTimeout: defaultHeaderTimeout,
	}
	return &Server{
		srv: &srv,
	}
}

func (s *Server) RegisterHandlers(ctx context.Context, pm *printer.PrinterManipulator) error {
	pingHandlers := handlers.InitBaseHandlers()
	printerHandlers := handlers.InitPrinterHandlers(pm)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		pingHandlers.Ping(w, r)
	})
	mux.HandleFunc("GET /api/v1/printer/find", func(w http.ResponseWriter, r *http.Request) {
        printerHandlers.FindPrinter(w, r)
    })
	mux.HandleFunc("POST /api/v1/printer/print", func(w http.ResponseWriter, r *http.Request) {
        printerHandlers.Print(w, r)
    })


	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	wrapper := logger.LoggerMidleware()(ctx, mux)
	s.srv.Handler = wrapper
	return nil
}


func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

