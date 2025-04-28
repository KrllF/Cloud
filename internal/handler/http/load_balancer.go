package http

import (
	"context"
	"encoding/json"
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
		h.Logg.Error("сервисы недоступны")
		body := ErrorResponce{
			Code:    http.StatusServiceUnavailable,
			Message: "сервисы недоступны",
		}
		b, err := json.Marshal(body)
		if err != nil {
			log.Printf("json.Marshal: %v\n", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := w.Write(b); err != nil {
			log.Printf("ошибка при отправке ответа: %v\n", err)
		}

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

				ctx := context.WithValue(request.Context(), consts.Attempts, attempts+1)

				h.LB(writer, request.WithContext(ctx))
			}

			h.proxies[peer.URL.Host] = proxy
		}
		h.mu.Unlock()
		proxy.ServeHTTP(w, r)

		return
	}
	h.Logg.Error("нет доступных серверов")
	body := ErrorResponce{
		Code:    http.StatusInternalServerError,
		Message: "нет доступных серверов",
	}
	b, err := json.Marshal(body)
	if err != nil {
		log.Printf("json.Marshal: %v\n", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write(b); err != nil {
		log.Printf("ошибка при отправке ответа: %v\n", err)
	}
}
