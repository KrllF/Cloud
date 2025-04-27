package middleware

import (
	"log"
	"net/http"
)

type rateLimiter interface {
	AddUser(id string) bool
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
func RateLimiter(TB rateLimiter) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getUserIP(r)
			ok := TB.AddUser(ip)
			if ok {
				log.Println("пользователь создан")
			} else {
				log.Println("пользователь существует")
			}

			_ = TB.Allow(ip)

			handler.ServeHTTP(w, r)
		})
	}
}
