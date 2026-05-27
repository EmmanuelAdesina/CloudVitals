package scorer

import "github.com/EmmanuelAdesina/CloudVitals/internal/core"

func Calculate(results []core.CheckResult) int {
	score := 100

	for _, r := range results {
		if r.Status != "fail" {
			continue
		}

		weight := 5
		switch r.Severity {
		case core.Critical:
			weight = 15
		case core.High:
			weight = 10
		case core.Medium:
			weight = 5
		case core.Low:
			weight = 2
		}

		score -= weight * len(r.Findings)
	}

	if score < 0 {
		return 0
	}
	return score
}
