package tools
import "math"

type bmiCalculator struct{
	weight float64
	height float64
}

func BmiCalculator(weight float64, height float64) float64 {
	calculatedBmi:= weight/(math.Pow(height,2))
	return calculatedBmi
}