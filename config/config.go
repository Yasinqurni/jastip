package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		DB
		AWS
		JWT
		Mongo
		GRPC `yaml:"grpc"`
		Upstore
	}

	// App -.
	App struct {
		Name     string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version  string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		TimeZone string `yaml:"timezone" env:"APP_TIMEZONE" `
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// JWT -.
	JWT struct {
		SECRET string `env-required:"true" env:"JWT_SECRET"`
		TTL    int64  `yaml:"jwt_ttl" env:"JWT_TTL" env-default:"60"`
	}

	// DB
	DB struct {
		URL string `env-required:"true" env:"DB_URL"`
	}

	// AWS
	AWS struct {
		AccessKeyID          string `env:"AWS_ACCESS_KEY_ID"`
		SecretAcessKey       string `env:"AWS_SECRET_ACCESS_KEY"`
		DefaultRegion        string `env:"AWS_DEFAULT_REGION"`
		Bucket               string `env:"AWS_BUCKET"`
		UsePathStyleEndpoint bool   `env:"AWS_USE_PATH_STYLE_ENDPOINT"`
	}

	// Mongo
	Mongo struct {
		DSN string `env-required:"true" env:"MONGO_DSN"`
		DB  string `env-required:"true" env:"MONGO_DB"`
	}

	// gRPC
	GRPC struct {
		Port string `env:"GRPC_PORT" yaml:"port"`
	}

	// Upstore
	Upstore struct {
		ApiBaseUrl string `env-required:"true" env:"UPSTORE_API_BASE_URL"`
		ApiToken   string `env-required:"true" env:"UPSTORE_API_BACKSTAGE_TOKEN"`
		WalletID   string `env-required:"true" env:"UPSTORE_WALLET_ID"`
		// e.g. /api/v8.0/...
		ApiVBaseUrl       string `env-required:"true" env:"UPSTORE_API_V_URL"`
		GoldpriceApiToken string `env-required:"true" env:"UPSTORE_GOLDPRICE_API_TOKEN"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Config error : ", err)
	}
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
