package calculator

import (
	"net/http"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/httputil"
)

func routes(h *Handler) []httputil.Route {
	return []httputil.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/calculate",
			Handler: h.calculate,
		},
	}
}
