package service

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/bool64/brick"
)

// Name is the name of this application or service.
const Name = "brick-template"

// Config defines application configuration.
type Config struct {
	brick.BaseConfig

	Jaeger jaeger.Options `split_words:"true"`
}
