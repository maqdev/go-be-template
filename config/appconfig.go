package config

import (
	"log/slog"
	"time"
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
	Host     string
	Username string
	Password string
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
