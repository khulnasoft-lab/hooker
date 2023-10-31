package webserver

import (
	"net/http"

	"github.com/khulnasoft-lab/hooker/v2/router"
)

func (web *WebServer) reload(w http.ResponseWriter, r *http.Request) {
	router.Instance().ReloadConfig()
}
