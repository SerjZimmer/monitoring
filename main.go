package main

import (
	"fmt"
	"github.com/SerjZimmer/monitoring/cmd/agent"
	"github.com/SerjZimmer/monitoring/cmd/server"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	go func() {
		// Создаем новый маршрутизатор Gorilla Mux
		r := mux.NewRouter()

		// Определяем маршрут для POST-запроса
		r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", server.UpdateHandler).Methods("POST")
		r.HandleFunc("/value/{metricType}/{metricName}", server.ValueHandler).Methods("GET")
		r.HandleFunc("/", server.ValueListHandler).Methods("GET")

		// Запускаем HTTP-сервер с использованием маршрутизатора Gorilla Mux
		http.Handle("/", r)
		fmt.Println("Сервер запущен на порту :8080")
		http.ListenAndServe(":8080", nil)
	}()
	go agent.Monitoring()

	select {}
}
