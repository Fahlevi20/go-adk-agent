package main

import ("fmt"
"github.com/Fahlevi20/go-adk-agent/tools"
)

func main() {
	bmi:=tools.BmiCalculator(70,1.75)
	fmt.Println("BMI Result:",bmi)
	fmt.Println("BMI Calculated")
}