package runner

import (
	"sync"

	"github.com/reliastra/cloudvitals/internal/core"
	"github.com/reliastra/cloudvitals/internal/providers"
)

func RunChecks(p providers.Provider, checks []core.CheckConfig, profile string, region string) []core.CheckResult {
	var wg sync.WaitGroup
	results := make([]core.CheckResult, len(checks))
	var mu sync.Mutex

	for i, check := range checks {
		wg.Add(1)
		go func(idx int, c core.CheckConfig) {
			defer wg.Done()

			res, err := p.RunCheck(c, profile, region)
			if err != nil {
				res = &core.CheckResult{
					CheckID:   c.ID,
					CheckName: c.Name,
					Status:    "error",
					Severity:  c.Severity,
					ErrorMsg:  err.Error(),
				}
			}

			mu.Lock()
			results[idx] = *res
			mu.Unlock()
		}(i, check)
	}

	wg.Wait()
	return results
}
