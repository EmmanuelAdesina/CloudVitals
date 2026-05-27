package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/reliastra/cloudvitals/internal/core"
	"github.com/reliastra/cloudvitals/internal/providers/aws"
	"github.com/reliastra/cloudvitals/internal/renderer"
	"github.com/reliastra/cloudvitals/internal/runner"
	"github.com/reliastra/cloudvitals/internal/scorer"
)

func main() {
	profile := flag.String("profile", "default", "AWS profile name")
	region := flag.String("region", "us-east-1", "AWS region")
	flag.Parse()

	registry, err := core.LoadRegistry("config/checks.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
		os.Exit(1)
	}

	provider := aws.New()
	checks := registry.ChecksForProvider("aws")

	if len(checks) == 0 {
		fmt.Println("No checks configured for provider: aws")
		os.Exit(0)
	}

	results := runner.RunChecks(provider, checks, *profile, *region)
	score := scorer.Calculate(results)
	renderer.Render(results, score)
}
