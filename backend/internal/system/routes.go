package system

import (
	"net/http"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/httputil"
)

func routes(h *Handler) []httputil.Route {
	return []httputil.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: h.health},
	}
}
