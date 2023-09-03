package main

import (
	"fmt"
	"github.com/SerjZimmer/monitoring/cmd/agent"
	"github.com/SerjZimmer/monitoring/cmd/server"
	"net/http"
)

func main() {
	go func() {
		http.HandleFunc("/update/", server.UpdateHandler)
		fmt.Println("Сервер запущен на порту 8080...")
		http.ListenAndServe(":8080", nil)
	}()
	go agent.Monitoring()

	select {}
}
