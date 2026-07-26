package proxy

import (
	"net/http"
	"strings"
)

func isWebSocket(r *http.Request) bool {

	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") &&
		strings.EqualFold(r.Header.Get("Connection"), "Upgrade")
}
