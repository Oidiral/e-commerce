package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Env      string `mapstructure:"ENV"`
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         string        `mapstructure:"SERVER_PORT"`
	ReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	User         string        `mapstructure:"DB_USER"`
	Password     string        `mapstructure:"DB_PASSWORD"`
	Host         string        `mapstructure:"DB_HOST"`
	Port         int           `mapstructure:"DB_PORT"`
	Name         string        `mapstructure:"DB_NAME"`
	SSLMode      string        `mapstructure:"DB_SSLMODE"`
	MaxOpenConns int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnTimeout  time.Duration `mapstructure:"DB_CONN_TIMEOUT"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"JWT_SECRET"`
	AccessTokenTTL  time.Duration `mapstructure:"JWT_ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `mapstructure:"JWT_REFRESH_TOKEN_TTL"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.SetEnvPrefix("AUTH")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}
