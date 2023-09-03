package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	// Создаем тестовый HTTP запрос
	req, err := http.NewRequest("GET", "/update/metricType/metricName/metricValue", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый HTTP ResponseWriter
	rr := httptest.NewRecorder()

	// Вызываем ваш обработчик
	UpdateHandler(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем содержимое ответа
	expectedResponse := "Метрика успешно принята: metricType/metricName/metricValue\n"
	assert.Equal(t, expectedResponse, rr.Body.String())
}

func TestHandler(t *testing.T) {
	// Создаем тестовый HTTP запрос
	req, err := http.NewRequest("GET", "/update/metricType/metricName/metricValue", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый HTTP ResponseWriter
	rr := httptest.NewRecorder()

	// Вызываем ваш обработчик
	http.HandlerFunc(UpdateHandler).ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateHandlerBadRequest(t *testing.T) {
	// Create a test HTTP request with an invalid URL
	req, err := http.NewRequest("GET", "/update/invalid/url", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test HTTP ResponseWriter
	rr := httptest.NewRecorder()

	// Call your updateHandler
	UpdateHandler(rr, req)

	// Check if the response status code is http.StatusBadRequest
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check the response body for the error message
	expectedResponse := "Неверный формат URL\n"
	assert.Equal(t, expectedResponse, rr.Body.String())
}
