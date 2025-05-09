package http

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/KrllF/Cloud/entity"
	"go.uber.org/zap"
)

const (
	sizeMap = 3
)

type (
	// ServerPool интерфейс пула бекенд серверов
	ServerPool interface {
		GetNextPeer() *entity.Backend
		MarkBackendStatus(backendURL *url.URL, alive bool)
	}
	// RateLimiter интерфейс rate-limiter
	RateLimiter interface {
		UpdateUser(ctx context.Context, ip string, tokenSize int64) error
	}
	// Handler структура хендлера
	Handler struct {
		Logg        *zap.Logger
		ServerPool  ServerPool
		RateLimiter RateLimiter
		proxies     map[string]*httputil.ReverseProxy
		mu          sync.RWMutex
	}
)

// NewHandler конструктор хендлера
func NewHandler(logg *zap.Logger, serverPool ServerPool, rateLimiter RateLimiter) Handler {
	return Handler{
		Logg:       logg,
		ServerPool: serverPool, RateLimiter: rateLimiter,
		proxies: make(map[string]*httputil.ReverseProxy, sizeMap), mu: sync.RWMutex{},
	}
}

// Init инициализация хендлера с middleware
func (h *Handler) Init(rateLimiterMiddleware func(http.Handler) http.Handler,
	logMiddleware func(http.Handler) http.Handler,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux := http.NewServeMux()

		mux.Handle("/", logMiddleware(rateLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.LB(w, r)
		}))))

		mux.HandleFunc("/client", h.UpdateTokenSize)

		mux.ServeHTTP(w, r)
	})
}
