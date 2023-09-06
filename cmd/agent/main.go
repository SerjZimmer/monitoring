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

var MetricsMap = map[string]float64{
	"Alloc":         0,
	"BuckHashSys":   0,
	"Frees":         0,
	"GCCPUFraction": 0,
	"GCSys":         0,
	"HeapAlloc":     0,
	"HeapIdle":      0,
	"HeapInuse":     0,
	"HeapObjects":   0,
	"HeapReleased":  0,
	"HeapSys":       0,
	"LastGC":        0,
	"Lookups":       0,
	"MCacheInuse":   0,
	"MCacheSys":     0,
	"MSpanInuse":    0,
	"MSpanSys":      0,
	"Mallocs":       0,
	"NextGC":        0,
	"NumForcedGC":   0,
	"NumGC":         0,
	"OtherSys":      0,
	"PauseTotalNs":  0,
	"StackInuse":    0,
	"StackSys":      0,
	"Sys":           0,
	"TotalAlloc":    0,
	"PollCount":     0,
	"RandomValue":   0,
}

func Monitoring() {
	// Создаем экземпляр структуры Metrics

	var mu sync.Mutex
	// Горутина для сбора метрик
	go func() {
		for {
			// Собираем метрики о памяти
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			// Обновляем метрики из пакета runtime
			mu.Lock()
			MetricsMap["Alloc"] = float64(m.Alloc)
			MetricsMap["BuckHashSys"] = float64(m.BuckHashSys)
			MetricsMap["Frees"] = float64(m.Frees)
			MetricsMap["GCCPUFraction"] = m.GCCPUFraction
			MetricsMap["GCSys"] = float64(m.GCSys)
			MetricsMap["HeapAlloc"] = float64(m.HeapAlloc)
			MetricsMap["HeapIdle"] = float64(m.HeapIdle)
			MetricsMap["HeapInuse"] = float64(m.HeapInuse)
			MetricsMap["HeapObjects"] = float64(m.HeapObjects)
			MetricsMap["HeapReleased"] = float64(m.HeapReleased)
			MetricsMap["HeapSys"] = float64(m.HeapSys)
			MetricsMap["LastGC"] = float64(m.LastGC)
			MetricsMap["Lookups"] = float64(m.Lookups)
			MetricsMap["MCacheInuse"] = float64(m.MCacheInuse)
			MetricsMap["MCacheSys"] = float64(m.MCacheSys)
			MetricsMap["MSpanInuse"] = float64(m.MSpanInuse)
			MetricsMap["MSpanSys"] = float64(m.MSpanSys)
			MetricsMap["Mallocs"] = float64(m.Mallocs)
			MetricsMap["NextGC"] = float64(m.NextGC)
			MetricsMap["NumForcedGC"] = float64(m.NumForcedGC)
			MetricsMap["NumGC"] = float64(m.NumGC)
			MetricsMap["OtherSys"] = float64(m.OtherSys)
			MetricsMap["PauseTotalNs"] = float64(m.PauseTotalNs)
			MetricsMap["StackInuse"] = float64(m.StackInuse)
			MetricsMap["StackSys"] = float64(m.StackSys)
			MetricsMap["Sys"] = float64(m.Sys)
			MetricsMap["TotalAlloc"] = float64(m.TotalAlloc)
			MetricsMap["PollCount"] = MetricsMap["PollCount"] + 1
			MetricsMap["RandomValue"] = rand.Float64()
			mu.Unlock()

			time.Sleep(pollInterval)
		}
	}()

	// Горутина для отправки метрик на сервер
	go func() {
		for {

			for metricName, metricValue := range MetricsMap {
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
}
