package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	sc "github.com/dmitry-konchenko/go-first-floor-sprint-four/internal/spentcalories"
)

// Основные ошибки, проверяемые программой
var (
	ErrInvalidFormat    = errors.New("Неверный формат строки: ожидается 'кол-во шагов, продолжительность'")
	ErrInvalidSteps     = errors.New("Неверное значение шагов: ожидается целочисленное значение")
	ErrNegativeSteps    = errors.New("Количество шагов должно быть больше 0")
	ErrInvalidDuration  = errors.New("Неверный формат продолжительности: ожидается формат типа 3h50m")
	ErrNegativeDuration = errors.New("Продолжительность должна быть больше 0")
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage преобразует строку с информацией о прогулке в кол-во шагов и продолжительность прогулки
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, ErrInvalidFormat
	}
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, ErrInvalidSteps
	}
	if steps <= 0 {
		return 0, 0, ErrNegativeSteps
	}
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, ErrInvalidDuration
	}
	if duration <= 0 {
		return 0, 0, ErrNegativeDuration
	}
	return steps, duration, nil
}

// DayActionInfo возвращает информацию о дневной активности: колличество шагов, дистанция, сожженные килокалории
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка парсинга программы", err)
		return ""
	}
	if steps <= 0 {
		log.Println("Ошибка колличества шагов", ErrNegativeSteps)
		return ""
	}
	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / mInKm
	calories, err := sc.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка вычисления затраченных калорий", err)
		return ""
	}
	answer := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
	return answer
}
