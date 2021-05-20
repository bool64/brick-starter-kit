package infra

import (
	"github.com/bool64/brick"
	"github.com/bool64/brick-template/internal/domain/greeting"
	"github.com/bool64/brick-template/internal/infra/schema"
	"github.com/bool64/brick-template/internal/infra/service"
	"github.com/bool64/brick/jaeger"
	"github.com/swaggest/rest/response/gzip"
)

// NewServiceLocator creates application service locator.
func NewServiceLocator(cfg service.Config) (*service.Locator, error) {
	bl, err := brick.NewBaseLocator(cfg.BaseConfig)
	if err != nil {
		return nil, err
	}

	l := &service.Locator{BaseLocator: bl}

	if err := jaeger.Setup(cfg.Jaeger, cfg.ServiceName, bl); err != nil {
		return nil, err
	}

	schema.SetupOpenapiCollector(l.OpenAPI)

	l.HTTPServerMiddlewares = append(l.HTTPServerMiddlewares, gzip.Middleware)

	l.GreetingMakerProvider = &greeting.SimpleMaker{}

	return l, nil
}
