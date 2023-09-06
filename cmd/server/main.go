package server

import (
	"fmt"
	"github.com/SerjZimmer/monitoring/cmd/agent"
	"net/http"
	"sort"
	"strings"
	"sync"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	// Разбираем URL-параметры
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	metricType := parts[2]
	metricName := parts[3]
	metricValue := parts[4]

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "Неверный тип метрики", http.StatusBadRequest)
		return
	}

	if _, ok := agent.MetricsMap[metricName]; !ok {
		http.Error(w, "Неверное имя метрики", http.StatusBadRequest)
		return
	}

	// Возвращаем успешный ответ
	fmt.Fprintf(w, "Метрика успешно принята: %s/%s/%s\n", metricType, metricName, metricValue)
}

func ValueHandler(w http.ResponseWriter, r *http.Request) {

	var mu sync.Mutex
	var result float64

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	// Разбираем URL-параметры
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	metricType := parts[2]
	metricName := parts[3]

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "Неверный тип метрики", http.StatusNotFound)
		return
	}

	if _, ok := agent.MetricsMap[metricName]; !ok {
		http.Error(w, "Неверное имя метрики", http.StatusNotFound)
		return
	}
	mu.Lock()
	result = agent.MetricsMap[metricName]
	mu.Unlock()

	// Возвращаем успешный ответ
	fmt.Fprintf(w, "%s/%s/ = %v\n", metricType, metricName, result)
}

func ValueListHandler(w http.ResponseWriter, r *http.Request) {
	// Заголовок ответа, указываем тип контента как text/html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Извлекаем ключи из мапы MetricsMap
	var keys []string
	for key := range agent.MetricsMap {
		keys = append(keys, key)
	}

	// Сортируем ключи
	sort.Strings(keys)

	// Генерируем HTML страницу
	fmt.Fprintf(w, "<html><head><title>Metrics</title></head><body>")
	fmt.Fprintf(w, "<h1>Все метрики</h1>")
	fmt.Fprintf(w, "<ul>")
	for _, key := range keys {
		fmt.Fprintf(w, "<li>%s: %v</li>", key, agent.MetricsMap[key])
	}
	fmt.Fprintf(w, "</ul></body></html>")
}
