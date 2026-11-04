package system

import (
	"net/http"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/httputil"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) error {
	return httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
