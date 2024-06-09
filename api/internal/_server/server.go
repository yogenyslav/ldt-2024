package server

import (
	"context"
	"fmt"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Nerzal/gocloak/v13"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/config"
	ac "github.com/yogenyslav/ldt-2024/api/internal/api/auth/controller"
	ah "github.com/yogenyslav/ldt-2024/api/internal/api/auth/handler"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg/metrics"
	"github.com/yogenyslav/ldt-2024/api/third_party"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server main struct that holds dependencies.
type Server struct {
	cfg      *config.Config
	srv      *grpc.Server
	pg       storage.SQLDatabase
	exporter sdktrace.SpanExporter
	tracer   trace.Tracer
}

// New creates a new Server instance.
func New(cfg *config.Config) *Server {
	var grpcOpts []grpc.ServerOption
	srv := grpc.NewServer(grpcOpts...)

	exporter := tracing.MustNewExporter(context.Background(), cfg.Jaeger.URL())
	provider := tracing.MustNewTraceProvider(exporter, "api")
	otel.SetTracerProvider(provider)
	tracer := otel.Tracer("api")

	return &Server{
		cfg:      cfg,
		srv:      srv,
		pg:       postgres.MustNew(cfg.Postgres, tracer),
		exporter: exporter,
		tracer:   tracer,
	}
}

// Run setups the server and starts it.
func (s *Server) Run() {
	defer s.pg.Close()
	defer func() {
		if err := s.exporter.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown exporter")
		}
	}()
	defer s.srv.GracefulStop()

	m := metrics.New()
	m.Collect()

	authController := ac.New(gocloak.NewClient(s.cfg.KeyCloak.URL), s.cfg.KeyCloak, s.tracer)
	authHandler := ah.New(authController, s.tracer, m)
	pb.RegisterAuthServiceServer(s.srv, authHandler)

	log.Info().Msg("starting the server")
	go s.listen()
	go s.listenGateway()
	go prom.HandlePrometheus(s.cfg.Prometheus)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("shutting down the server")
}

func (s *Server) listen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Server.Port))
	if err != nil {
		log.Panic().Err(err).Msg("failed to listen")
	}

	if err = s.srv.Serve(lis); err != nil {
		log.Error().Err(err).Msg("failed to start the server")
	}
}

func (s *Server) listenGateway() {
	var grpcOpts []grpc.DialOption
	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", s.cfg.Server.Port), grpcOpts...)

	if err != nil {
		log.Panic().Err(err).Msg("failed to connect to grpc server")
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close grpc connection")
		}
	}()

	mux := runtime.NewServeMux()
	if err = pb.RegisterAuthServiceHandler(context.Background(), mux, conn); err != nil {
		log.Panic().Err(err).Msg("failed to register the auth gateway ah")
	}

	withCors := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"ACCEPT", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v1") {
			mux.ServeHTTP(w, r)
			return
		}
		getOpenAPIHandler().ServeHTTP(w, r)
	}))

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Server.GatewayPort),
		Handler: withCors,
	}

	if err = gwServer.ListenAndServe(); err != nil { //nolint:G114 // not a security issue
		log.Error().Err(err).Msg("failed to start the gateway server")
	}
}

func getOpenAPIHandler() http.Handler {
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
	if err != nil {
		log.Error().Err(err).Msg("couldn't create sub filesystem")
	}
	return http.FileServer(http.FS(subFS))
}
