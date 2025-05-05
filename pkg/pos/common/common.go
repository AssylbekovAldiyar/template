package common

import (
	"pos-rs/pkg/pos/model"
)

type Config struct {
	Port string
	Env  string
	DB   struct {
		DSN string
	}
}

type Application struct {
	Config Config
	Models model.Models
}
