package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud-disk/core/router"
	"cloud-disk/core/svc"
	cfg "cloud-disk/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var configPath = flag.String("config", "configs/config.yaml", "Path to the service config file")

func main() {
	flag.Parse()

	conf, err := cfg.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	fx.New(
		fx.Supply(conf),
		fx.Provide(
			newLogger,
			provideServiceContext,
			provideGinEngine,
		),
		fx.Invoke(
			registerLoggerSync,
			registerServiceShutdown,
			registerHTTPServer,
		),
	).Run()
}

func newLogger(cfg cfg.Config) (*zap.Logger, error) {
	if cfg.Log.Mode == "console" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

func provideServiceContext(cfg cfg.Config) (*svc.ServiceContext, error) {
	return svc.NewServiceContext(cfg)
}

func provideGinEngine(cfg cfg.Config, svcCtx *svc.ServiceContext) *gin.Engine {
	if cfg.Log.Mode != "console" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	// Expose test static assets under /test for quick manual verification.
	engine.Static("/test", "./test")
	router.Register(engine, cfg.Name, svcCtx)
	return engine
}

func registerLoggerSync(lc fx.Lifecycle, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return logger.Sync()
		},
	})
}

func registerServiceShutdown(lc fx.Lifecycle, svcCtx *svc.ServiceContext) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return svcCtx.Close(ctx)
		},
	})
}

type serverParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    cfg.Config
	Engine    *gin.Engine
	Logger    *zap.Logger
}

func registerHTTPServer(p serverParams) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", p.Config.Host, p.Config.Port),
		Handler: p.Engine,
	}
	if secs := p.Config.HTTP.ReadTimeoutSeconds; secs > 0 {
		server.ReadTimeout = time.Duration(secs) * time.Second
	}
	if secs := p.Config.HTTP.WriteTimeoutSeconds; secs > 0 {
		server.WriteTimeout = time.Duration(secs) * time.Second
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info("starting http server", zap.String("addr", server.Addr))
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					p.Logger.Error("http server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("stopping http server")
			return server.Shutdown(ctx)
		},
	})
}
