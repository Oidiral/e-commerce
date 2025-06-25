package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Env             string         `mapstructure:"cart_env"`
	Server          ServerConfig   `mapstructure:",squash"`
	Database        DatabaseConfig `mapstructure:",squash"`
	Redis           Redis          `mapstructure:",squash"`
	CatalogGRPCAddr string         `mapstructure:"catalog_grpc_addr"`
	AuthURL         string         `mapstructure:"auth_url"`
	ClientID        string         `mapstructure:"client_id"`
	ClientSecret    string         `mapstructure:"client_secret"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"cart_server_port"`
	ReadTimeout  time.Duration `mapstructure:"cart_server_read_timeout"`
	WriteTimeout time.Duration `mapstructure:"cart_server_write_timeout"`
}

type Redis struct {
	Host     string `mapstructure:"cart_redis_host"`
	Password string `mapstructure:"cart_redis_password" default:""`
}

type DatabaseConfig struct {
	User         string        `mapstructure:"cart_db_user"`
	Password     string        `mapstructure:"cart_db_password"`
	Host         string        `mapstructure:"cart_db_host"`
	Port         int           `mapstructure:"cart_db_port"`
	Name         string        `mapstructure:"cart_db_name"`
	SSLMode      string        `mapstructure:"cart_db_sslmode"`
	MaxOpenConns int           `mapstructure:"cart_db_max_open_conns"`
	MaxIdleConns int           `mapstructure:"cart_db_max_idle_conns"`
	ConnTimeout  time.Duration `mapstructure:"cart_db_conn_timeout"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("warning: .env file not found or cannot be read, relying on system environment variables")
	}

	if err := viper.BindEnv("cart_env", "CART_ENV"); err != nil {
		return nil, fmt.Errorf("BindEnv CART_ENV: %w", err)
	}

	_ = viper.BindEnv("cart_server_port", "CART_SERVER_PORT")
	_ = viper.BindEnv("cart_server_read_timeout", "CART_SERVER_READ_TIMEOUT")
	_ = viper.BindEnv("cart_server_write_timeout", "CART_SERVER_WRITE_TIMEOUT")

	_ = viper.BindEnv("cart_db_user", "CART_DB_USER")
	_ = viper.BindEnv("cart_db_password", "CART_DB_PASSWORD")
	_ = viper.BindEnv("cart_db_host", "CART_DB_HOST")
	_ = viper.BindEnv("cart_db_port", "CART_DB_PORT")
	_ = viper.BindEnv("cart_db_name", "CART_DB_NAME")
	_ = viper.BindEnv("cart_db_sslmode", "CART_DB_SSLMODE")
	_ = viper.BindEnv("cart_db_max_open_conns", "CART_DB_MAX_OPEN_CONNS")
	_ = viper.BindEnv("cart_db_max_idle_conns", "CART_DB_MAX_IDLE_CONNS")
	_ = viper.BindEnv("cart_db_conn_timeout", "CART_DB_CONN_TIMEOUT")
	_ = viper.BindEnv("cart_redis_host", "CART_REDIS_HOST")
	_ = viper.BindEnv("cart_redis_password", "CART_REDIS_PASSWORD")
	_ = viper.BindEnv("catalog_grpc_addr", "CATALOG_GRPC_ADDR")
	_ = viper.BindEnv("auth_url", "AUTH_URL")
	_ = viper.BindEnv("client_id", "CLIENT_ID")
	_ = viper.BindEnv("client_secret", "CLIENT_SECRET")
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("viper.Unmarshal: %w", err)
	}

	if cfg.Env == "" {
		return nil, fmt.Errorf("CART_ENV must be set (in .env or environment)")
	}
	if cfg.Database.MaxOpenConns < 1 || cfg.Database.MaxIdleConns < 1 {
		return nil, fmt.Errorf("CART_DB_MAX_OPEN_CONNS and CART_DB_MAX_IDLE_CONNS must be >= 1")
	}
	if cfg.AuthURL == "" || cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("AUTH_URL, CLIENT_ID and CLIENT_SECRET must be set")
	}

	return &cfg, nil
}
