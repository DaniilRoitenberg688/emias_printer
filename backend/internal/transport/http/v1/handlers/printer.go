package handlers

import (
	"emias_printer/pkg/logger"
	"emias_printer/pkg/printer"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type CheckPrinterRequest struct {
	Ip string `json:"ip"`
}

type CheckPrinterResponse struct {
	Ip        string `json:"ip"`
	Available bool   `json:"available"`
}

type PrintRequest struct {
	Ip   string `json:"ip"`
	Text string `json:"text"`
}

type PrintResponse struct {
	Result string `json:"result"`
}

type PrinterHandlers struct {
	pm *printer.PrinterManipulator
}

func InitPrinterHandlers(pm *printer.PrinterManipulator) *PrinterHandlers {
	return &PrinterHandlers{pm: pm}
}

// FindPrinter находит ip принтера
// @Summary Находит ip принтера
// @Description Возвращает ip принтера
// @Tags Printer
// @Success 200 {string} string "ip"
// @Router /api/v1/printer/find [get]
func (h *PrinterHandlers) FindPrinter(w http.ResponseWriter, r *http.Request) {
	ips, err := h.pm.Scan()
	if err == printer.NoPrinterFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string][]string{"ip": ips}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.GetLoggerFromContext(r.Context()).Warn(r.Context(), "cannot encode data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Print отправляет текст на принтер
// @Summary Отправляет текст на принтер
// @Description Отправляет текст на принтер
// @Tags Printer
// @Accept json
// @Produce json
// @Param request body PrintRequest true "Request body"
// @Success 200 {object} PrintResponse
// @Router /api/v1/printer/print [post]
func (h *PrinterHandlers) Print(w http.ResponseWriter, r *http.Request) {
	var d PrintRequest
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		logger.GetLoggerFromContext(r.Context()).Warn(r.Context(), "cannot decode data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(d)

	if err := h.pm.SendRequest(d.Text, d.Ip, 9100); err != nil {
		logger.GetLoggerFromContext(r.Context()).Warn(r.Context(), "cannot send request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PrintResponse{"ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.GetLoggerFromContext(r.Context()).Warn(r.Context(), "cannot encode data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Check Printer проверяет доступен ли принтер
// @Summary проверяет доступен ли принтер
// @Tags Printer
// @Accept json
// @Produce json
// @Param request body CheckPrinterRequest true "Request body"
// @Success 200 {object} CheckPrinterResponse
// @Router /api/v1/printer/check [post]
func (h *PrinterHandlers) Check(w http.ResponseWriter, r *http.Request) {
	var d CheckPrinterRequest
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		logger.GetLoggerFromContext(r.Context()).Warn(r.Context(), "cannot decode data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	av := h.pm.CheckPrinterIp(d.Ip)
	data := CheckPrinterResponse{
		Ip:        d.Ip,
		Available: av,
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.GetLoggerFromContext(r.Context()).Warn(r.Context(), "cannot encode data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
