package http

import (
	"encoding/json"
	"log"
	"net/http"

	"go.uber.org/zap"
)

// UpdateTokenSizeRequest запрос
type UpdateTokenSizeRequest struct {
	UserIP    string `json:"user_ip"`
	TokenSize int64  `json:"token_size"`
}

// UpdateTokenSize обновить максильманый размер токенов
func (h *Handler) UpdateTokenSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logg.Info("некорректный метод")
		body := ErrorResponce{
			Code:    http.StatusBadRequest,
			Message: "r.Method != http.MethodPost",
		}
		b, err := json.Marshal(body)
		if err != nil {
			log.Printf("json.Marshal: %v\n", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(b); err != nil {
			log.Printf("ошибка при отправке ответа: %v\n", err)
		}

		return
	}
	defer r.Body.Close()

	var token UpdateTokenSizeRequest
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		body := ErrorResponce{
			Code:    http.StatusBadRequest,
			Message: "json.NewDecoder",
		}
		b, err := json.Marshal(body)
		if err != nil {
			log.Printf("json.Marshal: %v\n", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(b); err != nil {
			log.Printf("ошибка при отправке ответа: %v\n", err)
		}

		return
	}

	if err = h.RateLimiter.UpdateUser(r.Context(), token.UserIP, token.TokenSize); err != nil {
		h.Logg.Error("h.RateLimiter.UpdateUser", zap.Error(err))
		body := ErrorResponce{
			Code:    http.StatusInternalServerError,
			Message: "h.RateLimiter.UpdateUser",
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

		return
	}

	w.WriteHeader(http.StatusOK)
}
