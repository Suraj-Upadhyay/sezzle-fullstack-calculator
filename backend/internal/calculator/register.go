package calculator

import (
	"net/http"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/httputil"
)

func Register(mux *http.ServeMux) {
	httputil.Register(mux, routes(NewHandler(NewService())))
}
