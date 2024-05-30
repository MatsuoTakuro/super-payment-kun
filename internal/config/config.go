package config

import "github.com/caarlos0/env/v6"

type Config struct {
	HttpPort   int    `env:"HTTP_PORT" envDefault:"8080"`
	DBHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"DB_PORT" envDefault:"3316"`
	DBUser     string `env:"DB_USER" envDefault:"taro"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"pass"`
	DBName     string `env:"DB_NAME" envDefault:"super-payment-kun-db"`
	JWTSecret  string `env:"JWT_SECRET" envDefault:"MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNA=="`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
