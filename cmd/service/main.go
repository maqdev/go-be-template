package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	api "github.com/maqdev/go-be-template/gen/api/authors"
	"github.com/maqdev/go-be-template/service"

	"github.com/maqdev/go-be-template/config"
	"github.com/maqdev/go-be-template/util/logutil"
)

func main() {
	if err := run(); err != nil {
		os.Exit(-1)
	}
}

func run() error {
	configPath := flag.String("config", "", "Path to the config file (required)")
	flag.Parse()
	if *configPath == "" {
		flag.Usage()
		return errors.New("No config file was provided")
	}

	const serviceName = "YEY"
	log := logutil.NewStdLogger(slog.LevelInfo)
	log.Info("Starting up service", "service", serviceName, "config", *configPath, "version", config.VersionString())

	cfg, err := config.LoadConfig(serviceName, *configPath)
	if err != nil {
		log.Error("LoadConfig failed", "err", err)
		return err
	}

	log = logutil.NewStdLogger(cfg.Log.Level.SLogLevel())
	slog.SetDefault(log)

	handler := service.NewHandler(cfg)
	server, err := api.NewServer(handler)
	if err != nil {
		log.Error("Server initialization failed", "err", err)
		return err
	}

	log.Info("Listening and serving", "address", cfg.HTTP.Address)

	httpServer := http.Server{
		Addr:                         cfg.HTTP.Address,
		Handler:                      server,
		DisableGeneralOptionsHandler: true,
		TLSConfig:                    nil,
		ReadTimeout:                  cfg.HTTP.ReadTimeout,
		WriteTimeout:                 cfg.HTTP.WriteTimeout,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer close(signalChan)
	go handleClose(log, cfg, &httpServer, signalChan)

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Couldn't listen and serve", "err", err)
		return err
	}

	return nil
}

func handleClose(log *slog.Logger, cfg *config.AppConfig, httpServer *http.Server, signalChan chan os.Signal) {
	sig := <-signalChan

	if cfg.HTTP.ShutdownDelay > 0 {
		log.Info("Shutdown signal was received, delaying", "signal", sig, "timeout", cfg.HTTP.ShutdownDelay)
		time.Sleep(cfg.HTTP.ShutdownDelay)
	}

	log.Info("Shutting down http server", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Error("Couldn't shutdown gracefully", "err", err)
	}
}
