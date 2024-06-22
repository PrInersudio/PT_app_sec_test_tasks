package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// структуры хранения конфигурации
type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	StopTimeout time.Duration `yaml:"stop_timeout" env-default:"1ms"`
	RateLimit   `yaml:"rate limit"`
}

type RateLimit struct {
	Limit    int           `yaml:"limit" env-default:"100"`
	Interval time.Duration `yaml:"interval" env-default:"60s"`
	Msg      string        `yaml:"msg" env-default:"Слишком много запросов."`
}

// загрузка конфигурации из файла
func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		log.Panicf("Ошибка открытия файла конфигурации %s", err)
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Panicf("Ошибка чтения файла конфигурации %s", err)
	}
	return &cfg
}
