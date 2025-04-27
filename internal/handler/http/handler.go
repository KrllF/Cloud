package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/KrllF/Cloud/entity"
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
	// Handler структура хендлера
	Handler struct {
		ServerPool ServerPool
		proxies    map[string]*httputil.ReverseProxy
		mu         sync.RWMutex
	}
)

// NewHandler конструктор хендлера
func NewHandler(serverPool ServerPool) Handler {
	return Handler{ServerPool: serverPool, proxies: make(map[string]*httputil.ReverseProxy, sizeMap), mu: sync.RWMutex{}}
}

// Init инициализация хендлера
func (h *Handler) Init() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.LB(w, r)
	})
}
