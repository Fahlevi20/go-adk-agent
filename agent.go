package main

import (
  "context"
  "log"
  "os"

  "google.golang.org/adk/agent/llmagent"
  "google.golang.org/adk/cmd/launcher/adk"
  "google.golang.org/adk/cmd/launcher/full"
  "google.golang.org/adk/model/gemini"
  "google.golang.org/adk/server/restapi/services"
  "google.golang.org/adk/tool"
  "google.golang.org/adk/tool/geminitool"
  "google.golang.org/genai"
  "github.com/Fahlevi20/go-adk-agent/tools"
)

func main() {
  ctx := context.Background()

  model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
  })
  if err != nil {
    log.Fatalf("Failed to create model: %v", err)
  }

  agent, err := llmagent.New(llmagent.Config{
    Name:        "bmi_agent",
    Model:       model,
    Description: "Tells the BMI (Body Mass Index) based on user weight and height.",
    Instruction: "You are a helpful assistant that tells the BMI based on user weight and height.",
    Tools: []tool.Tool{
      tools.bmiCalculator{},
    },
  })
  if err != nil {
    log.Fatalf("Failed to create agent: %v", err)
  }

  config := &adk.Config{
    AgentLoader: services.NewSingleAgentLoader(agent),
  }

  l := full.NewLauncher()
  err = l.Execute(ctx, config, os.Args[1:])
  if err != nil {
    log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
  }
}