package middleware

import (
	"context"
	"log"
	"net/http"
)

type rateLimiter interface {
	AddUser(ctx context.Context, id string) (bool, error)
	Allow(id string) bool
}

func getUserIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}

// RateLimiter middleware, которая проверяет token bucket для пользователя
func RateLimiter(tb rateLimiter) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getUserIP(r)
			ok, err := tb.AddUser(r.Context(), ip)
			if err != nil {
				http.Error(w, "ошибка на сервере", http.StatusInternalServerError)

				return
			}
			if ok {
				log.Println("пользователь создан")
			} else {
				log.Println("пользователь существует")
			}

			ok = tb.Allow(ip)
			if ok {
				handler.ServeHTTP(w, r)

				return
			}

			http.Error(w, "нет токенов", http.StatusTooManyRequests)
		})
	}
}
