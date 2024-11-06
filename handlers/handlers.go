package handlers

import (
	"backend/redis"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func sendError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// Получение значение переменной со временем последнего изменения (api/value)
func GetValueHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	value, lastModifiedTime, err := redis.GetValue()
	if err != nil {
		log.Printf("Ошибка при получении данных: %v", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при получении данных")
		return
	}

	response := map[string]string{
		"value":              value,
		"last_modified_time": lastModifiedTime,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при отправке данных: %v", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при отправке данных")
	}
}

// Обновляем значение переменной и изменяем время последнего изменение на текущее (api/value/update)
func UpdateValueHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Ошибка при разборе данных запроса: %v", err)
		sendError(w, http.StatusBadRequest, "Ошибка при разборе данных запроса")
		return
	}

	var newValue string
	if value, exists := requestBody["value"]; exists && value != "" {
		newValue = value
	} else {
		newValue = redis.GenerateRandomValue()
	}

	err := redis.SetValue(newValue)
	if err != nil {
		log.Printf("Ошибка при обновлении данных: %v", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при обновлении данных")
		return
	}

	response := map[string]string{
		"value":              newValue,
		"last_modified_time": time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при отправке данных: %v", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при отправке данных")
	}
}

// Оставляем значение переменной, но изменяем время последнего изменение на текущее (api/value/keep)
func KeepValueHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := redis.UpdateLastModifiedTime()
	if err != nil {
		log.Printf("Ошибка при обновлении времени: %v", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при обновлении времени")
		return
	}

	newLastModifiedTime := time.Now().Format(time.RFC3339)

	response := map[string]string{
		"last_modified_time": newLastModifiedTime,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при отправке времени: %v", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при отправке времени")
	}
}
