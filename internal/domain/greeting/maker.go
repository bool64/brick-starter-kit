package greeting

import (
	"context"
	"strings"

	"github.com/bool64/ctxd"
)

type Params struct {
	Name   string `query:"name" default:"World"`
	Locale string `query:"locale" required:"true" enum:"en-US,ru-RU"`
}

type Maker interface {
	Hello(ctx context.Context, params Params) (string, error)
}

type SimpleMaker struct{}

func (s *SimpleMaker) Hello(ctx context.Context, params Params) (string, error) {
	if strings.ToLower(params.Name) == "bug" {
		return "", ctxd.NewError(ctx, "#$@@^! %C 🤖")
	}

	switch params.Locale {
	case "en-US":
		return "Hello, " + params.Name + "!", nil
	case "ru-RU":
		return "Привет, " + params.Name + "!", nil
	default:
		return "", ctxd.NewError(ctx, "unknown locale", "locale", params.Locale)
	}
}

func (s *SimpleMaker) GreetingMaker() Maker {
	return s
}
