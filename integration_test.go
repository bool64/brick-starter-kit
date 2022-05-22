package main_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/bool64/brick"
	"github.com/bool64/brick-starter-kit/internal/infra"
	"github.com/bool64/brick-starter-kit/internal/infra/nethttp"
	"github.com/bool64/brick-starter-kit/internal/infra/service"
	"github.com/bool64/brick-starter-kit/internal/infra/storage"
	"github.com/bool64/brick/config"
	"github.com/bool64/brick/test"
	"github.com/bool64/httptestbench"
	"github.com/godogx/dbsteps"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestFeatures(t *testing.T) {
	var cfg service.Config

	test.RunFeatures(t, "", &cfg, func(tc *test.Context) (*brick.BaseLocator, http.Handler) {
		cfg.ServiceName = service.Name

		sl, err := infra.NewServiceLocator(cfg)
		require.NoError(t, err)

		tc.Database.Instances[dbsteps.Default] = dbsteps.Instance{
			Tables: map[string]interface{}{
				storage.GreetingsTable: new(storage.GreetingRow),
			},
		}

		return sl.BaseLocator, nethttp.NewRouter(sl)
	})
}

func BenchmarkGreetings(b *testing.B) {
	var cfg service.Config
	cfg.ServiceName = service.Name

	require.NoError(b, config.Load("", &cfg, config.WithOptionalEnvFiles(".env.integration-test")))

	sl, err := infra.NewServiceLocator(cfg)
	if err != nil {
		b.Skip(err)
	}

	router := nethttp.NewRouter(sl)

	srv := httptest.NewServer(router)
	defer srv.Close()

	httptestbench.RoundTrip(b, 50,
		func(i int, req *fasthttp.Request) {
			req.SetRequestURI(srv.URL + "/hello?locale=en-US&name=user" + strconv.Itoa(((i/10)^12345)%100))
		},
		func(i int, resp *fasthttp.Response) bool {
			return resp.StatusCode() == http.StatusOK
		},
	)
}
