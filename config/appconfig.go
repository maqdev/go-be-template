package config

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppConfig struct {
	Log  LogConfig
	HTTP HTTPConfig
	DB   DBConfig
}

type HTTPConfig struct {
	Address         string
	ShutdownDelay   time.Duration
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

type DBConfig struct {
	URL            string
	Password       string
	MaxConnections int32
}

type LogConfig struct {
	Level Level
}

type Level int

func (l *Level) UnmarshalString(s string) error {
	var slevel slog.Level
	if err := slevel.UnmarshalText([]byte(s)); err != nil {
		return err
	}
	*l = Level(slevel)
	return nil
}

func (l *Level) SLogLevel() slog.Level {
	return slog.Level(*l)
}

func (dbc DBConfig) CreatePool(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbc.URL)
	if err != nil {
		return nil, err
	}
	if dbc.Password != "" {
		config.ConnConfig.Password = dbc.Password
	}
	config.MaxConns = dbc.MaxConnections

	return pgxpool.NewWithConfig(ctx, config)
}
