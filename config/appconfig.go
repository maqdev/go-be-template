package config

type DBConfig struct {
	Host     string
	Username string
	Password string
}

type LogConfig struct {
}

type AppConfig struct {
	DB DBConfig
}
