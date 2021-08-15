package main

import (
	"fmt"
	"ova-checklist-api/internal/config"
)

func main() {
	simpleConfig := config.OpenSimpleConfig("configs/simple_config.json")
	fmt.Printf("Hi there! You are running ova-checklist-api. Current environment: %s\n", simpleConfig.Environment)
}
