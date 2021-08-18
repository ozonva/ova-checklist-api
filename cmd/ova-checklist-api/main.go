package main

import (
	"fmt"
	"ova-checklist-api/internal/config"
	"time"
)

func main() {
	fmt.Printf("Hi there! You are running ova-checklist-api\n")
	configReader := config.NewSimpleDynamicConfigReader(
		"configs/simple_dynamic_config.json",
		5 * time.Second,
	)
	configReader.Update()
}
