package redis

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
}

// Получаем значение и время последнего обновления
func GetValue() (string, string, error) {

	// Получаем значение
	value, err := rdb.Get(ctx, "random_value").Result()
	if err == redis.Nil {
		// Если значение не существует, генерируем новое
		value = GenerateRandomValue()
		rdb.Set(ctx, "random_value", value, 0)
	} else if err != nil {
		log.Printf("Ошибка при получении значения из Redis: %v", err)
		return "", "", err
	}

	// Получаем время последнего изменения
	lastModifiedTime, err := rdb.Get(ctx, "last_modified_time").Result()
	if err == redis.Nil {
		// Если времени нет, установим текущее время
		lastModifiedTime = time.Now().Format(time.RFC3339)
		rdb.Set(ctx, "last_modified_time", lastModifiedTime, 0)
	} else if err != nil {
		log.Printf("Ошибка при получении времени из Redis: %v", err)
		return "", "", err
	}

	return value, lastModifiedTime, nil
}

// Сохранение нового значения в Redis
func SetValue(newValue string) error {

	// Сохраняем новое значение в Redis
	err := rdb.Set(ctx, "random_value", newValue, 0).Err()
	if err != nil {
		log.Printf("Ошибка при сохранении значения в Redis: %v", err)
		return err
	}

	// Обновляем время последнего изменения
	currentTime := time.Now().Format(time.RFC3339)
	err = rdb.Set(ctx, "last_modified_time", currentTime, 0).Err()
	if err != nil {
		log.Printf("Ошибка при обновлении времени в Redis: %v", err)
	}
	return err
}

// Обновление времени последнего изменения
func UpdateLastModifiedTime() error {
	currentTime := time.Now().Format(time.RFC3339)
	return rdb.Set(ctx, "last_modified_time", currentTime, 0).Err()
}

// Генерация случайного значения
func GenerateRandomValue() string {
	return fmt.Sprintf("%d", rand.Intn(1000))
}
