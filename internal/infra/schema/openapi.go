package schema

import (
	"github.com/swaggest/rest/openapi"
)

func SetupOpenapiCollector(c *openapi.Collector) {
	c.Reflector().SpecEns().Info.Title = "brick-template"
}
