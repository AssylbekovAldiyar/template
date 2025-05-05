package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	HTTP       HttpConfig
	Logger     LoggerConfig
	Kafka      KafkaConfig
	HTTPClient HTTPClientConfig
}

type HttpConfig struct {
	Addr                  string `env:"HTTP_CONFIG_ADDR" default:":8080"`
	HealthADDR            string `env:"HTTP_CONFIG_HEALTH_ADDR" default:":8081"`
	RequestTimeoutSeconds int    `env:"HTTP_CONFIG__REQUEST_TIMEOUT_SECONDS"   default:"60"`
}

type LoggerConfig struct {
	Level  int    `env:"LOGGER_CONFIG_LEVEL" `
	Format string `env:"LOGGER_CONFIG_FORMAT" default:"json"`
}

type KafkaConfig struct {
	Brokers         []string `env:"KAFKA_BROKERS" default:"127.0.0.1:9092"`
	ConsumerGroup   string   `env:"KAFKA_CONSUMER_GROUP"`
	TimeOutSeconds  int      `env:"KAFKA_CONFIG_TIMEOUT_SECONDS" default:"60"`
	PartitionNumber int      `env:"KAFKA_CONFIG_PARTITION_NUMBER" default:"12"`
}

type HTTPClientConfig struct {
	MaxIdleConnections    int `json:"max_idle_connections"`
	MaxConnsPerHost       int `json:"max_conns_per_host"`
	IdleConnTimeout       int `json:"idle_conn_timeout"`
	ResponseHeaderTimeout int `json:"response_header_timeout"`
	Timeout               int `json:"timeout"` // in milliseconds
}

func New(filenames ...string) (*Config, error) {
	cfg := new(Config)

	if len(filenames) > 0 {
		if err := godotenv.Load(filenames...); err != nil {
			return nil, err
		}
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	defaults.SetDefaults(cfg)

	return cfg, nil
}
