package function_calling

import (
	"fmt"

	"google.golang.org/adk/tool" // Add this import
)

// BmiCalculatorArgs defines the input parameters
type BmiCalculatorArgs struct {
	Weight float64 `json:"weight" jsonschema:"Weight in kilograms"`
	Height float64 `json:"height" jsonschema:"Height in meters"`
}

// BmiCalculatorResult defines the output schema for the bmi_calculator tool.
type BmiCalculatorResult struct {
	BMI      float64 `json:"bmi"`
	Category string  `json:"category"`
	Message  string  `json:"message"`
}

// CalculateBMI calculates BMI based on weight and height
func CalculateBMI(ctx tool.Context, input BmiCalculatorArgs) BmiCalculatorResult {
	if input.Height <= 0 {
		return BmiCalculatorResult{
			BMI:      0,
			Category: "Error",
			Message:  "Height must be greater than 0",
		}
	}
	if input.Weight <= 0 {
		return BmiCalculatorResult{
			BMI:      0,
			Category: "Error", 
			Message:  "Weight must be greater than 0",
		}
	}

	calculatedBmi := input.Weight / (input.Height * input.Height)
	category := getBMICategory(calculatedBmi)
	
	fmt.Printf("Tool: Calculated BMI for weight %.1f kg and height %.2f m: %.1f (%s)\n", 
		input.Weight, input.Height, calculatedBmi, category)
		
	return BmiCalculatorResult{
		BMI:      calculatedBmi,
		Category: category,
		Message:  fmt.Sprintf("The BMI is %.1f (%s)", calculatedBmi, category),
	}
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