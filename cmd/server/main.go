package server

import (
	"fmt"
	"net/http"
	"strings"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Разбираем URL-параметры
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	metricType := parts[2]
	metricName := parts[3]
	metricValue := parts[4]

	// Возвращаем успешный ответ
	fmt.Fprintf(w, "Метрика успешно принята: %s/%s/%s\n", metricType, metricName, metricValue)
}
