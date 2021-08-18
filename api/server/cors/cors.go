package cors

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// CombinedCORSHandler wraps a server and provides CORS headers
func CombinedCORSHandler(h http.Handler) http.Handler {
	return corsHandler{h}
}

type corsHandler struct {
	handler http.Handler
}

func (c corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)

	if r.Method == "OPTIONS" {
		return
	}

	c.handler.ServeHTTP(w, r)
}

// SetHeaders sets the CORS headers
func SetHeaders(w http.ResponseWriter, r *http.Request) {
	set := func(w http.ResponseWriter, k, v string) {
		if v := w.Header().Get(k); len(v) > 0 {
			return
		}
		w.Header().Set(k, v)
	}

	if origin := r.Header.Get("Origin"); len(origin) > 0 {
		set(w, "Access-Control-Allow-Origin", origin)
	} else {
		set(w, "Access-Control-Allow-Origin", "*")
	}

	set(w, "Access-Control-Allow-Credentials", "true")
	set(w, "Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
	if v := os.Getenv("CORS_ALLOW_ORIGIN_HEADER"); strings.TrimSpace(v) != "" {
		set(w, "Access-Control-Allow-Headers", fmt.Sprintf("%s, %s", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization", v))
	} else {
		set(w, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
}
