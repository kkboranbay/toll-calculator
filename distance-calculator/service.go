package main

import (
	"math"

	"github.com/kkboranbay/toll-calculator/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculateService struct {
	prevPoints []float64
}

func NewCalculateService() CalculatorServicer {
	return &CalculateService{}
}

func (s *CalculateService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(s.prevPoints) > 0 {
		distance = calculateDistance(s.prevPoints[0], s.prevPoints[1], data.Lat, data.Long)
	}

	s.prevPoints = []float64{data.Lat, data.Long}

	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
