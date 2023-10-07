package main

import (
	"math"

	"github.com/kkboranbay/toll-calculator/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculateService struct {
	points [][]float64
}

func NewCalculateService() *CalculateService {
	return &CalculateService{
		points: make([][]float64, 0),
	}
}

func (s *CalculateService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(s.points) > 0 {
		prevPoint := s.points[len(s.points)-1]
		distance = calculateDistance(prevPoint[0], prevPoint[1], data.Lat, data.Long)
	}

	s.points = append(s.points, []float64{data.Lat, data.Long})

	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
