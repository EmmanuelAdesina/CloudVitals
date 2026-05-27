package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/EmmanuelAdesina/CloudVitals/internal/core"
	"github.com/EmmanuelAdesina/CloudVitals/internal/providers/aws"
	"github.com/EmmanuelAdesina/CloudVitals/internal/providers/mock"
	"github.com/EmmanuelAdesina/CloudVitals/internal/renderer"
	"github.com/EmmanuelAdesina/CloudVitals/internal/runner"
	"github.com/EmmanuelAdesina/CloudVitals/internal/scorer"
)

func main() {
	profile := flag.String("profile", "default", "AWS profile name")
	region := flag.String("region", "us-east-1", "AWS region")
	mockMode := flag.Bool("mock", false, "Use mock provider for testing")
	flag.Parse()

	registry, err := core.LoadRegistry("config/checks.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
		os.Exit(1)
	}

	var provider core.Provider
	if *mockMode {
		provider = mock.New()
	} else {
		provider = aws.New()
	}

	checks := registry.ChecksForProvider("aws")
	if len(checks) == 0 {
		fmt.Println("No checks configured for provider: aws")
		os.Exit(0)
	}

	results := runner.RunChecks(provider, checks, *profile, *region)
	score := scorer.Calculate(results)
	renderer.Render(results, score)
}
