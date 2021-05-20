package ui

import (
	"net/http"
	"os"

	"github.com/bool64/brick-template/resources/static"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

var Static http.Handler

func init() {
	if _, err := os.Stat("./resources/static"); err == nil {
		// path/to/whatever exists
		Static = http.FileServer(http.Dir("./resources/static"))
	} else {
		Static = statigz.FileServer(static.Assets, brotli.AddEncoding, statigz.EncodeOnInit)
	}
}

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Static.ServeHTTP(w, r)
	})
}
