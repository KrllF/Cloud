package http

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/KrllF/Cloud/internal/consts"
)

const (
	countAttempts       = 3
	countRetry          = 3
	timeWaitMillisecond = 10
)

// LB балансировщик запросов
func (h *Handler) LB(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > countAttempts {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)

		return
	}

	peer := h.ServerPool.GetNextPeer()
	if peer != nil {
		h.mu.Lock()
		proxy, ok := h.proxies[peer.URL.Host]
		if !ok {
			proxy = httputil.NewSingleHostReverseProxy(peer.URL)
			proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
				log.Printf("[%s] %s\n", peer.URL.Host, e.Error())

				retries := GetRetryFromContext(request)
				if retries < countRetry {
					<-time.After(timeWaitMillisecond * time.Millisecond)
					ctx := context.WithValue(request.Context(), consts.Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))

					return
				}

				h.ServerPool.MarkBackendStatus(peer.URL, false)

				attempts := GetAttemptsFromContext(request)
				log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
				ctx := context.WithValue(request.Context(), consts.Attempts, attempts+1)

				h.LB(writer, request.WithContext(ctx))
			}

			h.proxies[peer.URL.Host] = proxy
		}
		h.mu.Unlock()
		proxy.ServeHTTP(w, r)

		return
	}

	http.Error(w, "нет доступных серверов", http.StatusInternalServerError)
}
