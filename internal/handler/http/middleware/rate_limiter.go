package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type rateLimiter interface {
	AddUser(ctx context.Context, id string) (bool, error)
	Allow(id string) bool
}

// getUserIP получить ip пользователя
func getUserIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip != "" && isValidIP(ip) {
		return ip
	}

	forwared := r.Header.Get("X-Forwarded-For")
	if forwared != "" {
		ips := strings.Split(forwared, ",")
		clientIP := strings.TrimSpace(ips[0])
		if isValidIP(clientIP) {
			return clientIP
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && isValidIP(host) {
		return host
	}

	return ""
}

// isValidIP проверяет является ли валидным ip
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// RateLimiter middleware, которая проверяет token bucket для пользователя
func RateLimiter(logg *zap.Logger, tb rateLimiter) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getUserIP(r)
			ok, err := tb.AddUser(r.Context(), ip)
			if err != nil {
				logg.Error("ошибка на сервере", zap.Error(err))
				http.Error(w, "ошибка на сервере", http.StatusInternalServerError)

				return
			}
			if ok {
				logg.Info("пользователь создан")
			}
			ok = tb.Allow(ip)
			if ok {
				handler.ServeHTTP(w, r)

				return
			}
			logg.Error("нет токенов", zap.String("ip", ip))
			http.Error(w, "нет токенов", http.StatusTooManyRequests)
		})
	}
}
