package config

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-faster/errors"
	"github.com/redis/go-redis/v9"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppConfig struct {
	Log   LogConfig
	HTTP  HTTPConfig
	DB    DBConfig
	Redis RedisConfig
}

type LogConfig struct {
	Level Level
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

type RedisConfig struct {
	Opts        *redis.UniversalOptions
	ClusterMode bool
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

func (rc RedisConfig) CreateClient() (redis.UniversalClient, error) {
	if rc.Opts == nil {
		return nil, errors.New("RedisConfig.Opts is nil")
	}
	if rc.ClusterMode {
		return redis.NewClusterClient(rc.Opts.Cluster()), nil
	} else {
		return redis.NewUniversalClient(rc.Opts), nil
	}
}
