package function_calling

import (
	"context"
	"fmt"
	"math"

	"google.golang.org/adk/tool"
)

type BmiCalculator struct{}

func (b BmiCalculator) Name() string {
	return "bmi_calculator"
}

func (b BmiCalculator) Description() string {
	return "Calculates BMI (Body Mass Index) based on weight (kg) and height (m)."
}

func (b BmiCalculator) Schema() tool.Schema {
	return tool.Schema{
		Type: "object",
		Properties: map[string]tool.Property{
			"weight": {
				Type:        "number",
				Description: "Weight in kilograms",
			},
			"height": {
				Type:        "number", 
				Description: "Height in meters",
			},
		},
		Required: []string{"weight", "height"},
	}
}

func (b BmiCalculator) Call(ctx context.Context, input map[string]interface{}) (*tool.Result, error) {
	weight, ok := input["weight"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'weight' parameter")
	}
	
	height, ok := input["height"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'height' parameter")
	}

	if height <= 0 {
		return nil, fmt.Errorf("height must be greater than 0")
	}

	calculatedBmi := weight / (math.Pow(height, 2))
	
	// Determine BMI category
	category := b.getBMICategory(calculatedBmi)
	
	result := fmt.Sprintf("The BMI is %.2f (%s)", calculatedBmi, category)
	
	return &tool.Result{
		Data: result,
	}, nil
}

func (b BmiCalculator) getBMICategory(bmi float64) string {
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