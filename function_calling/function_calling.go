package function_calling

import (
	"context"
	"fmt"
)

// BmiCalculatorArgs defines the input parameters
type BmiCalculatorArgs struct {
	Weight float64 `json:"weight" jsonschema:"Weight in kilograms"`
	Height float64 `json:"height" jsonschema:"Height in meters"`
}

// BmiCalculatorResult defines the output structure
type BmiCalculatorResult struct {
	BMI      float64 `json:"bmi"`
	Category string  `json:"category"`
	Message  string  `json:"message"`
}

// CalculateBMI must return (BmiCalculatorResult, error)
func CalculateBMI(ctx context.Context, args BmiCalculatorArgs) (BmiCalculatorResult, error) {
	if args.Height <= 0 {
		return BmiCalculatorResult{
			BMI:      0,
			Category: "Error",
			Message:  "Height must be greater than 0",
		}, fmt.Errorf("height must be greater than 0")
	}

	calculatedBmi := args.Weight / (args.Height * args.Height)
	category := getBMICategory(calculatedBmi)
	
	return BmiCalculatorResult{
		BMI:      calculatedBmi,
		Category: category,
		Message:  fmt.Sprintf("The BMI is %.2f (%s)", calculatedBmi, category),
	}, nil
}

func getBMICategory(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "Underweight"
	case bmi < 25:
		return "Normal weight"
	case bmi < 30:
		return "Overweight"
	default:
		return "Obese"
	}
}