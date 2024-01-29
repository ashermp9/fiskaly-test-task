package ports

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ashermp9/fiskaly-test-task/internal/app"
	"github.com/ashermp9/fiskaly-test-task/internal/ports/types"
	"go.uber.org/zap"
)

type Server struct {
	logger        *zap.SugaredLogger
	APIService    *app.APIService
	server        *http.Server
	listenAddress int
}

func NewServer(logger *zap.SugaredLogger, appService *app.APIService, listenAddress int) *Server {
	return &Server{
		APIService:    appService,
		listenAddress: listenAddress,
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", listenAddress),
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/create-device", s.LoggingMiddleware(http.HandlerFunc(s.CreateSignatureDeviceHandler)))
	mux.Handle("/api/v0/sign-transaction", s.LoggingMiddleware(http.HandlerFunc(s.SignTransactionHandler)))
	mux.Handle("/api/v0/health", s.LoggingMiddleware(http.HandlerFunc(s.HealthCheckHandler)))

	s.server.Handler = mux // Set the mux as the server's handler
	return s.server.ListenAndServe()
}

func (s *Server) SignTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var signRequest types.SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&signRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	domainRequest := types.ConvertToDomainSignTransactionRequest(signRequest)
	signature, err := s.APIService.SignTransaction(ctx, domainRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(signature)
}

func (s *Server) CreateSignatureDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var createDeviceRequest types.CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&createDeviceRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := createDeviceRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	domainRequest := types.ConvertToDomainCreateDeviceRequest(createDeviceRequest)
	device, err := s.APIService.CreateDevice(ctx, domainRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(device)
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	healthStatus := map[string]string{"status": "pass", "version": "v0"}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthStatus)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// LoggingMiddleware logs each request's method, URL, and duration.
func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		s.logger.Infof("Handled request: %s %s, Duration: %s", r.Method, r.URL.Path, duration)
	})
}
