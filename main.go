package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений.
const (
	MInKm      = 1000
	MinInHours = 60
	LenStep    = 0.65
	CmInM      = 100
)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string
	Action       int
	LenStep      float64
	Duration     time.Duration
	Weight       float64
}

// distance возвращает дистанцию, которую преодолел пользователь.
// Формула расчета:
// количество_повторов * длина_шага / м_в_км
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	return t.distance() / (t.Duration.Hours())
}

// Calories возвращает количество потраченных килокалорий на тренировке.
// Пока возвращаем 0, так как этот метод будет переопределяться для каждого типа тренировки.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// TrainingInfo возвращает труктуру InfoMessage, в которой хранится вся информация о проведенной тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String возвращает строку с информацией о проведенной тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для структур: Running, Walking и Swimming.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для расчета потраченных килокалорий при беге.
const (
	CaloriesMeanSpeedMultiplier = 18
	CaloriesMeanSpeedShift      = 1.79
)

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training
}

// Calories возввращает количество потраченных килокалория при беге.
// Формула расчета:
// ((18 * средняя_скорость_в_км/ч + 1.79) * вес_спортсмена_в_кг / м_в_км * время_тренировки_в_часах * мин_в_часе)
// Это переопределенный метод Calories() из Training.
func (r Running) Calories() float64 {
	return ((CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Hours() * MinInHours)
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training.
func (r Running) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: r.TrainingType,
		Duration:     r.Duration,
		Distance:     r.distance(),
		Speed:        r.meanSpeed(),
		Calories:     r.Calories(),
	}
}

// Константы для расчета потраченных килокалорий при ходьбе.
const (
	CaloriesWeightMultiplier      = 0.035
	CaloriesSpeedHeightMultiplier = 0.029
	KmHInMsec                     = 0.278
)

// Walking структура описывающая тренировку Ходьба
type Walking struct {
	Training
	Height float64
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
// Формула расчета:
// ((0.035 * вес_спортсмена_в_кг + (средняя_скорость_в_метрах_в_секунду**2 / рост_в_метрах)
// * 0.029 * вес_спортсмена_в_кг) * время_тренировки_в_часах * мин_в_ч)
// Это переопределенный метод Calories() из Training.
func (w Walking) Calories() float64 {
	speedInMsec := w.meanSpeed() * KmHInMsec
	return ((CaloriesWeightMultiplier*w.Weight + (math.Pow(speedInMsec, 2)/(w.Height/CmInM))*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours)
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training.
func (w Walking) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: w.TrainingType,
		Duration:     w.Duration,
		Distance:     w.distance(),
		Speed:        w.meanSpeed(),
		Calories:     w.Calories(),
	}
}

// Константы для расчета потраченных килокалорий при плавании.
const (
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2
)

// Swimming структура, описывающая тренировку Плавание
type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

// meanSpeed возвращает среднюю скорость при плавании.
// Формула расчета:
// длина_бассейна * количество_пересечений / м_в_км / продолжительность_тренировки
// Это переопределенный метод Calories() из Training.
func (s Swimming) meanSpeed() float64 {
	return float64(s.LengthPool*s.CountPool) / MInKm / s.Duration.Hours()
}

// Calories возвращает количество калорий, потраченных при плавании.
// Формула расчета:
// (средняя_скорость_в_км/ч + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * вес_спортсмена_в_кг * время_тренировки_в_часах
// Это переопределенный метод Calories() из Training.
func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

// TrainingInfo returns info about swimming training.
// Это переопределенный метод TrainingInfo() из Training.
func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     float64(s.LengthPool*s.CountPool) / MInKm,
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

// ReadData возвращает информацию о проведенной тренировке.
func ReadData(training CaloriesCalculator) string {
	calories := training.Calories()
	info := training.TrainingInfo()
	return fmt.Sprintf("%s\nПотрачено ккал: %.2f\n", info.String(), calories)
}

func main() {

	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
