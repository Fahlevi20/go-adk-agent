package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
		
	"github.com/joho/godotenv" // Add this import for .env support

)

// BmiCalculatorArgs defines the schema for the arguments passed to the bmi_calculator tool.
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

// CalculateBMI is a tool that calculates BMI based on weight and height.
// Note: Uses tool.Context instead of context.Context like in the reference example
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

func main() {
		// Check if API key is set
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding to check environment variables")
	}
	ctx := context.Background()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable is not set")
	}

	model, err := gemini.NewModel(ctx, "gemini-2.0-flash-exp", &genai.ClientConfig{
		APIKey: apiKey,
	})
	
	fmt.Println("✅ Model created successfully")
	fmt.Println("✅ API Key is present")

	if err != nil {
		log.Fatalf("❌ Failed to create model: %v", err)
	}
	
	fmt.Println("🔄 Creating agent...")
	
	// Create BMI tool using functiontool.New() - exactly like the reference
	bmiTool, err := functiontool.New(
		functiontool.Config{
			Name:        "bmi_calculator",
			Description: "Calculates BMI (Body Mass Index) based on weight (kg) and height (m).",
		},
		CalculateBMI, // Our function that uses tool.Context
	)
	if err != nil {
		log.Fatalf("❌ Failed to create BMI tool: %v", err)
	}
	
	agent, err := llmagent.New(llmagent.Config{
		Name:        "bmi_calculator",
		Model:       model,
		Description: "Calculates BMI (Body Mass Index) based on user weight and height.",
		Instruction: `You are a helpful BMI calculator assistant. 
When users provide their weight and height, use the bmi_calculator tool to compute their BMI.
Always provide the BMI value along with the corresponding category (Underweight, Normal weight, Overweight, or Obese).
Be friendly and encouraging in your responses.`,
		Tools: []tool.Tool{
			bmiTool,
		},
	})
	
	if err != nil {
		log.Fatalf("❌ Failed to create agent: %v", err)
	}
	
	fmt.Println("✅ Agent setup complete")

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}

	l := full.NewLauncher()
	
	fmt.Println("🚀 Starting agent server...")
	fmt.Println("The agent will be available for BMI calculations")
	fmt.Println("You can interact with it using the appropriate client")
	
	err = l.Execute(ctx, config, os.Args[1:])
	if err != nil {
		log.Fatalf("❌ Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}