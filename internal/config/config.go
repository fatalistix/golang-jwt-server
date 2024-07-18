package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string                `yaml:"env" env-required:"true"`
	Gin             GinConfig             `yaml:"gin" env-required:"true"`
	Postgres        PostgresConfig        `yaml:"postgres" env-required:"true"`
	Token           TokenConfig           `yaml:"token" env-required:"true"`
	Http            HttpConfig            `yaml:"http" env-required:"true"`
	PasswordEncoder PasswordEncoderConfig `yaml:"password-encoder" env-required:"true"`
}

type GinConfig struct {
	RunMode string `yaml:"runmode" env-required:"true"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Db       string `env:"POSTGRES_DB" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	SslMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
	Port     uint16 `env:"POSTGRES_PORT" env-required:"true"`
}

type TokenConfig struct {
	Secret          string        `env:"SECRET" env-required:"true"`
	AccessTokenTtl  time.Duration `yaml:"access-token-ttl" env-required:"true"`
	RefreshTokenTtl time.Duration `yaml:"refresh-token-ttl" env-required:"true"`
	Issuer          string        `yaml:"issuer" env-required:"true"`
	CleanTimeout    time.Duration `yaml:"clean-timeout" env-required:"true"`
}

type HttpConfig struct {
	Port            uint16        `yaml:"port" env-required:"true"`
	ReadTimeout     time.Duration `yaml:"read-timeout" env-required:"true"`
	WriteTimeout    time.Duration `yaml:"write-timeout" env-required:"true"`
	IdleTimeout     time.Duration `yaml:"idle-timeout" env-required:"true"`
	ShutdownTimeout time.Duration `yaml:"shutdown-timeout" env-required:"true"`
}

type PasswordEncoderConfig struct {
	EncryptCost int `env:"ENCRYPT_COST" env-required:"true"`
}

func MustLoad(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists: " + configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return config
}
