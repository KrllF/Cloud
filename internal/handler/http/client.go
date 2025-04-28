package http

import (
	"encoding/json"
	"net/http"
)

// UpdateTokenSizeRequest запрос
type UpdateTokenSizeRequest struct {
	UserIP    string `json:"user_ip"`
	TokenSize int64  `json:"token_size"`
}

// UpdateTokenSize обновить максильманый размер токенов
func (h *Handler) UpdateTokenSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	var token UpdateTokenSizeRequest
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = h.RateLimiter.UpdateUser(r.Context(), token.UserIP, token.TokenSize); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
