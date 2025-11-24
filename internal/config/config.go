package config

import (
    "os"
)

type Config struct {
    ServerAddress string `env:"SERVER_ADDRESS"`
    Port          string `env:"SERVER_PORT"`
    DbUser        string `env:"DB_USER"`
    DbPassword    string `env:"DB_PASSWORD"`
    DbHost        string `env:"DB_HOST"`
    DbPort        string `env:"DB_PORT"`
    SslMode       string `env:"SSL_MODE"`
    DbName        string `env:"DB_NAME"`
    JWTSecret     string `env:"JWT_SECRET"`
}

func NewConfig() Config {
    return Config{}
}

func (cfg *Config) InitConfig() error {
    cfg.ServerAddress = os.Getenv("SERVER_ADDRESS")
    cfg.Port = os.Getenv("SERVER_PORT")
    cfg.DbUser = os.Getenv("DB_USER")
    cfg.DbPassword = os.Getenv("DB_PASSWORD")
    cfg.DbHost = os.Getenv("DB_HOST")
    cfg.DbPort = os.Getenv("DB_PORT")
    cfg.SslMode = os.Getenv("SSL_MODE")
    cfg.DbName = os.Getenv("DB_NAME")
    return nil
}