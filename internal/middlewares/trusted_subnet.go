package middlewares

import (
	"net"
	"net/http"

	"github.com/SpaceSlow/execenv/internal/config"
)

// WithCheckingTrustedSubnet middleware предназначена для проверки исходящих запросов из доверенной подсети.
func WithCheckingTrustedSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, err := config.GetServerConfig()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		realIP := r.Header.Get("X-Real-IP")
		if cfg.TrustedSubnet != config.NewCIDR("") && !cfg.TrustedSubnet.Contains(net.ParseIP(realIP)) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
