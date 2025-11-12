package tools
import "math"

type bmiCalculator struct{
	weight float64
	height float64
}

func (s bmiCalculator) Name() string(
	return "bmi_calculator"
)

func (s bmiCalculator) Description() string{
	return "Calculates BMI (Body Mass Index) based on weight (kg) and height (m)."
}
func BmiCalculator(weight float64, height float64) float64 {
	calculatedBmi:= weight/(math.Pow(height,2))
	return calculatedBmi
}