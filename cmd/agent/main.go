package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func Monitoring() {
	// Создаем экземпляр структуры Metrics
	metrics := make(map[string]float64)
	var mu sync.Mutex
	var PollCount float64
	// Горутина для сбора метрик
	go func() {
		for {
			// Собираем метрики о памяти
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			// Обновляем метрики из пакета runtime
			mu.Lock()
			metrics["Alloc"] = float64(m.Alloc)
			metrics["BuckHashSys"] = float64(m.BuckHashSys)
			metrics["Frees"] = float64(m.Frees)
			metrics["GCCPUFraction"] = m.GCCPUFraction
			metrics["GCSys"] = float64(m.GCSys)
			metrics["HeapAlloc"] = float64(m.HeapAlloc)
			metrics["HeapIdle"] = float64(m.HeapIdle)
			metrics["HeapInuse"] = float64(m.HeapInuse)
			metrics["HeapObjects"] = float64(m.HeapObjects)
			metrics["HeapReleased"] = float64(m.HeapReleased)
			metrics["HeapSys"] = float64(m.HeapSys)
			metrics["LastGC"] = float64(m.LastGC)
			metrics["Lookups"] = float64(m.Lookups)
			metrics["MCacheInuse"] = float64(m.MCacheInuse)
			metrics["MCacheSys"] = float64(m.MCacheSys)
			metrics["MSpanInuse"] = float64(m.MSpanInuse)
			metrics["MSpanSys"] = float64(m.MSpanSys)
			metrics["Mallocs"] = float64(m.Mallocs)
			metrics["NextGC"] = float64(m.NextGC)
			metrics["NumForcedGC"] = float64(m.NumForcedGC)
			metrics["NumGC"] = float64(m.NumGC)
			metrics["OtherSys"] = float64(m.OtherSys)
			metrics["PauseTotalNs"] = float64(m.PauseTotalNs)
			metrics["StackInuse"] = float64(m.StackInuse)
			metrics["StackSys"] = float64(m.StackSys)
			metrics["Sys"] = float64(m.Sys)
			metrics["TotalAlloc"] = float64(m.TotalAlloc)
			metrics["PollCount"] = PollCount + 1
			metrics["RandomValue"] = rand.Float64()
			mu.Unlock()

			time.Sleep(pollInterval)
		}
	}()

	// Горутина для отправки метрик на сервер
	go func() {
		for {

			for metricName, metricValue := range metrics {
				if metricName != "PollCount" {
					go sendMetric("gauge", metricName, metricValue)
				} else {
					go sendMetric("counter", metricName, int64(metricValue))
				}
			}

			time.Sleep(reportInterval)
		}
	}()

	// Ожидание завершения работы агента
	select {}
}

// Отправка метрики на сервер
func sendMetric(metricType, metricName string, metricValue any) {
	serverURL := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metricType, metricName, metricValue)

	// Отправляем POST-запрос на сервер
	req, err := http.NewRequest("POST", serverURL, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	req.Header.Set("Content-Type", "text/plain")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке метрики на сервер:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка при отправке метрики на сервер. Код ответа:", resp.StatusCode)
		return
	}

	fmt.Printf("Метрика успешно отправлена на сервер: %s/%s/%f\n", metricType, metricName, metricValue)
}
