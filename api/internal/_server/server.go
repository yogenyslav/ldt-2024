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
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/config"
	ac "github.com/yogenyslav/ldt-2024/api/internal/api/auth/controller"
	ah "github.com/yogenyslav/ldt-2024/api/internal/api/auth/handler"
	authmw "github.com/yogenyslav/ldt-2024/api/internal/api/auth/middleware"
	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/repo"
	"github.com/yogenyslav/ldt-2024/api/internal/api/middleware"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/internal/api/predictor/handler"
	ph "github.com/yogenyslav/ldt-2024/api/internal/api/prompter/handler"
	sc "github.com/yogenyslav/ldt-2024/api/internal/api/stock/controller"
	sh "github.com/yogenyslav/ldt-2024/api/internal/api/stock/handler"
	sr "github.com/yogenyslav/ldt-2024/api/internal/api/stock/repo"
	"github.com/yogenyslav/ldt-2024/api/pkg/client"
	"github.com/yogenyslav/ldt-2024/api/pkg/metrics"
	"github.com/yogenyslav/ldt-2024/api/third_party"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage"
	"github.com/yogenyslav/pkg/storage/mongo"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server струткура сервера со всеми зависимостями.
type Server struct {
	cfg      *config.Config
	srv      *grpc.Server
	pg       storage.SQLDatabase
	mongo    storage.MongoDatabase
	kc       *gocloak.GoCloak
	exporter sdktrace.SpanExporter
	tracer   trace.Tracer
}

// New создает новый Server.
func New(cfg *config.Config) *Server {
	kc := gocloak.NewClient(cfg.KeyCloak.URL)

	exporter := tracing.MustNewExporter(context.Background(), cfg.Jaeger.URL())
	provider := tracing.MustNewTraceProvider(exporter, "api")
	otel.SetTracerProvider(provider)
	tracer := otel.Tracer("api")

	pg := postgres.MustNew(cfg.Postgres, tracer)

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}
	var grpcOpts []grpc.ServerOption
	grpcOpts = append(grpcOpts, grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(middleware.InterceptorLogger(), logOpts...),
		auth.UnaryServerInterceptor(authmw.JWT(kc, cfg.KeyCloak.Realm, repo.New(pg))),
	))
	srv := grpc.NewServer(grpcOpts...)

	return &Server{
		cfg:      cfg,
		srv:      srv,
		pg:       pg,
		mongo:    mongo.MustNew(cfg.Mongo, tracer),
		kc:       kc,
		exporter: exporter,
		tracer:   tracer,
	}
}

// Run создает все контроллеры и хендлеры и запускает сервер.
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

	authController := ac.New(s.kc, s.cfg.KeyCloak, s.tracer)
	authHandler := ah.New(authController, s.tracer, m)
	pb.RegisterAuthServiceServer(s.srv, authHandler)

	prompterClient, err := client.NewGrpcClient(s.cfg.Prompter)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create prompter grpc client")
	}
	defer func() {
		if err := prompterClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close prompter grpc client")
		}
	}()
	predictorClient, err := client.NewGrpcClient(s.cfg.Predictor)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create predictor grpc client")
	}
	defer func() {
		if err := prompterClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close predictor grpc client")
		}
	}()

	pb.RegisterPrompterServer(s.srv, ph.New(prompterClient, s.tracer))

	stockRepo := sr.New(s.mongo)
	stockController := sc.New(stockRepo, s.tracer)
	stockHandler := sh.New(stockController, s.tracer)
	pb.RegisterStockServer(s.srv, stockHandler)

	pb.RegisterPredictorServer(s.srv, handler.New(predictorClient, s.tracer))

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

	ctx := context.Background()
	mux := runtime.NewServeMux()
	if err = pb.RegisterAuthServiceHandler(ctx, mux, conn); err != nil {
		log.Panic().Err(err).Msg("failed to register the auth gateway")
	}
	if err = pb.RegisterPrompterHandler(ctx, mux, conn); err != nil {
		log.Panic().Err(err).Msg("failed to register the prompter gateway")
	}
	if err = pb.RegisterStockHandler(ctx, mux, conn); err != nil {
		log.Panic().Err(err).Msg("failed to register the stock gateway")
	}
	if err = pb.RegisterPredictorHandler(ctx, mux, conn); err != nil {
		log.Panic().Err(err).Msg("failed to register the predictor gateway")
	}

	withCors := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"ACCEPT", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin", "Origin", "Accept", "ngrok-skip-browser-warning"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("path", r.URL.Path).Msg("incoming request")
		if strings.HasPrefix(r.URL.Path, "/api/v1") {
			mux.ServeHTTP(w, r)
			return
		}
		getOpenAPIHandler().ServeHTTP(w, r)
	}))

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Server.GatewayPort),
		Handler: withCors,
	}

	if err = gwServer.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("failed to start the gateway server")
	}
}

func getOpenAPIHandler() http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		log.Error().Err(err).Msg("couldn't add mime type")
	}
	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
	if err != nil {
		log.Error().Err(err).Msg("couldn't create sub filesystem")
	}
	return http.FileServer(http.FS(subFS))
}
