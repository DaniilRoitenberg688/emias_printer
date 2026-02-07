package handlers

import (
	"encoding/json"
	"net/http"
)
type Handlers struct{}

func InitBaseHandlers() *Handlers {
	return &Handlers{}
}


// Ping проверяет доступность сервера
// @Summary Проверка доступности сервера
// @Description Возвращает "pong" если сервер работает
// @Tags Ping
// @Success 200 {string} string "pong"
// @Router /api/v1/ping [get]
func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"result": "ok"}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}


