package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные ошибки, проверяемые программой
var (
	ErrInvalidFormat    = errors.New("Неверный формат строки: ожидается 'кол-во шагов,активность,продолжительность'")
	ErrInvalidSteps     = errors.New("Неверное значение шагов: ожидается целочисленное значение")
	ErrNegativeSteps    = errors.New("Количество шагов не может быть отрицательным")
	ErrInvalidDuration  = errors.New("Неверный формат продолжительности: ожидается формат типа 3h50m")
	ErrNegativeDuration = errors.New("Продолжительность должна быть больше 0")
	ErrNegativeWeight   = errors.New("Вес должен быть больше 0")
	ErrNegativeHeight   = errors.New("Рост должен быть больше 0")
	ErrInvalidTraining  = errors.New("неизвестный тип тренировки")
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит строку формата "шаги,вид активности,продолжительность"
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, ErrInvalidFormat
	}
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, ErrInvalidSteps
	}
	if steps <= 0 {
		return 0, "", 0, ErrNegativeSteps
	}
	activity := parts[1]
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, ErrInvalidDuration
	}
	if duration <= 0 {
		return 0, "", 0, ErrNegativeDuration
	}
	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в километрах на основе количества шагов и роста
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm
	return distanceKm
}

// meanSpeed вычисляет среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	hours := duration.Hours()
	speed := distance / hours
	return speed
}

// TrainingInfo возвращает информацию о тренировке на основе входных данных
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	var distanceVal float64
	var speedVal float64
	var caloriesVal float64
	var errCal error
	switch activity {
	case "Ходьба":
		distanceVal = distance(steps, height)
		speedVal = meanSpeed(steps, height, duration)
		caloriesVal, errCal = WalkingSpentCalories(steps, weight, height, duration)
		if errCal != nil {
			return "", errCal
		}

	case "Бег":
		distanceVal = distance(steps, height)
		speedVal = meanSpeed(steps, height, duration)
		caloriesVal, errCal = RunningSpentCalories(steps, weight, height, duration)
		if errCal != nil {
			return "", errCal
		}

	default:
		return "", ErrInvalidTraining
	}
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), distanceVal, speedVal, caloriesVal)

	return result, nil
}

// RunningSpentCalories вычисляет количество калорий, потраченных при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrNegativeSteps
	}
	if weight <= 0 {
		return 0, ErrNegativeWeight
	}
	if height <= 0 {
		return 0, ErrNegativeHeight
	}
	if duration <= 0 {
		return 0, ErrNegativeDuration
	}
	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

// WalkingSpentCalories вычисляет количество калорий, потраченных при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrNegativeSteps
	}
	if weight <= 0 {
		return 0, ErrNegativeWeight
	}
	if height <= 0 {
		return 0, ErrNegativeHeight
	}
	if duration <= 0 {
		return 0, ErrNegativeDuration
	}
	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	baseCalories := (weight * speed * durationMinutes) / minInH
	calories := baseCalories * walkingCaloriesCoefficient
	return calories, nil
}
