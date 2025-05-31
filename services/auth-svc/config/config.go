package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Env      string         `mapstructure:"auth_env"`
	Server   ServerConfig   `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	JWT      JWTConfig      `mapstructure:",squash"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"auth_server_port"`
	ReadTimeout  time.Duration `mapstructure:"auth_server_read_timeout"`
	WriteTimeout time.Duration `mapstructure:"auth_server_write_timeout"`
}

type DatabaseConfig struct {
	User         string        `mapstructure:"auth_db_user"`
	Password     string        `mapstructure:"auth_db_password"`
	Host         string        `mapstructure:"auth_db_host"`
	Port         int           `mapstructure:"auth_db_port"`
	Name         string        `mapstructure:"auth_db_name"`
	SSLMode      string        `mapstructure:"auth_db_sslmode"`
	MaxOpenConns int           `mapstructure:"auth_db_max_open_conns"`
	MaxIdleConns int           `mapstructure:"auth_db_max_idle_conns"`
	ConnTimeout  time.Duration `mapstructure:"auth_db_conn_timeout"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"auth_jwt_secret"`
	AccessTokenTTL  time.Duration `mapstructure:"auth_jwt_access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"auth_jwt_refresh_token_ttl"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("warning: .env file not found or cannot be read, relying on system environment variables")
	}

	if err := viper.BindEnv("auth_env", "AUTH_ENV"); err != nil {
		return nil, fmt.Errorf("BindEnv AUTH_ENV: %w", err)
	}

	_ = viper.BindEnv("auth_server_port", "AUTH_SERVER_PORT")
	_ = viper.BindEnv("auth_server_read_timeout", "AUTH_SERVER_READ_TIMEOUT")
	_ = viper.BindEnv("auth_server_write_timeout", "AUTH_SERVER_WRITE_TIMEOUT")

	_ = viper.BindEnv("auth_db_user", "AUTH_DB_USER")
	_ = viper.BindEnv("auth_db_password", "AUTH_DB_PASSWORD")
	_ = viper.BindEnv("auth_db_host", "AUTH_DB_HOST")
	_ = viper.BindEnv("auth_db_port", "AUTH_DB_PORT")
	_ = viper.BindEnv("auth_db_name", "AUTH_DB_NAME")
	_ = viper.BindEnv("auth_db_sslmode", "AUTH_DB_SSLMODE")
	_ = viper.BindEnv("auth_db_max_open_conns", "AUTH_DB_MAX_OPEN_CONNS")
	_ = viper.BindEnv("auth_db_max_idle_conns", "AUTH_DB_MAX_IDLE_CONNS")
	_ = viper.BindEnv("auth_db_conn_timeout", "AUTH_DB_CONN_TIMEOUT")

	_ = viper.BindEnv("auth_jwt_secret", "AUTH_JWT_SECRET")
	_ = viper.BindEnv("auth_jwt_access_token_ttl", "AUTH_JWT_ACCESS_TOKEN_TTL")
	_ = viper.BindEnv("auth_jwt_refresh_token_ttl", "AUTH_JWT_REFRESH_TOKEN_TTL")

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("viper.Unmarshal: %w", err)
	}

	if cfg.Env == "" {
		return nil, fmt.Errorf("AUTH_ENV must be set (in .env or environment)")
	}
	if cfg.Database.MaxOpenConns < 1 || cfg.Database.MaxIdleConns < 1 {
		return nil, fmt.Errorf("AUTH_DB_MAX_OPEN_CONNS and AUTH_DB_MAX_IDLE_CONNS must be >= 1")
	}
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("AUTH_JWT_SECRET must be set")
	}

	return &cfg, nil
}
