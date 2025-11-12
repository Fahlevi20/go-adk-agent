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
	
	"github.com/Fahlevi20/go-adk-agent/function_calling" // Use your local module name
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding to check environment variables")
	}
	
	ctx := context.Background()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable is not set")
	}

	model, err := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{
		APIKey: apiKey,
	})
	
	fmt.Println("✅ Model created successfully")
	fmt.Println("✅ API Key is present")

	if err != nil {
		log.Fatalf("❌ Failed to create model: %v", err)
	}
	
	fmt.Println("🔄 Creating agent...")
	
	// Create BMI tool using functiontool.New()
	bmiTool, err := functiontool.New(
		functiontool.Config{
			Name:        "bmi_calculator",
			Description: "Calculates BMI (Body Mass Index) based on weight (kg) and height (m).",
		},
		function_calling.CalculateBMI,
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