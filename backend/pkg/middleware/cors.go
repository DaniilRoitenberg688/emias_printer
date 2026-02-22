package middleware

import (
	"slices"
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig defines the allowed origins, methods, headers, etc.
type CORSConfig struct {
	AllowedOrigins   []string // e.g. []string{"*"} or explicit origins
	AllowedMethods   []string // e.g. []string{"GET", "POST", "OPTIONS"}
	AllowedHeaders   []string // e.g. []string{"Content-Type", "Authorization"}
	ExposedHeaders   []string // optional headers the browser may read
	AllowCredentials bool
	MaxAge           int // seconds browsers may cache the pre‑flight response
}

// DefaultCORS returns a permissive config useful for development.
func DefaultCORS() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400,
	}
}

// CORSMiddleware returns an http.Handler that injects CORS headers and handles pre‑flight OPTIONS requests.
func CORSMiddleware(next http.Handler, cfg *CORSConfig) http.Handler {
	// Pre‑compute comma‑separated strings for performance.
	allowMethods := strings.Join(cfg.AllowedMethods, ", ")
	allowHeaders := strings.Join(cfg.AllowedHeaders, ", ")
	exposeHeaders := strings.Join(cfg.ExposedHeaders, ", ")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request – just continue.
			next.ServeHTTP(w, r)
			return
		}

		// Set Access‑Control‑Allow‑Origin
		if cfg.AllowedOrigins[0] == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// Echo back the origin only if it is whitelisted.
			if slices.Contains(cfg.AllowedOrigins, origin) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
		}
		w.Header().Add("Vary", "Origin")
		if cfg.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if exposeHeaders != "" {
			w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
		}

		// Pre‑flight request handling.
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			if cfg.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Normal request – continue down the chain.
		next.ServeHTTP(w, r)
	})
}
