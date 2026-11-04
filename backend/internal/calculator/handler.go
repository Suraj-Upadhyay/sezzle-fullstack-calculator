package calculator

import (
	"net/http"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/httputil"
)

type Calculator interface {
	Calculate(c Calculation) (float64, error)
}

type Handler struct {
	calculator Calculator
}

func NewHandler(calculator Calculator) *Handler {
	return &Handler{calculator: calculator}
}

func (h *Handler) calculate(w http.ResponseWriter, r *http.Request) error {
	req, err := httputil.DecodeJSON[CalculateRequest](w, r)
	if err != nil {
		return err
	}

	calculation := req.toDomain()
	result, err := h.calculator.Calculate(calculation)
	if err != nil {
		return toAPIError(err)
	}

	return httputil.WriteJSON(w, http.StatusOK,
		newCalculateResponse(calculation, result))
}
