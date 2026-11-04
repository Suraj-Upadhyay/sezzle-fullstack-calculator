package middleware

import "net/http"

// Wrap applies the global middleware. Recover sits inside Logger so a
// recovered panic is still logged with its 500 status.
func Wrap(h http.Handler) http.Handler {
	return Logger(Recover(h))
}
