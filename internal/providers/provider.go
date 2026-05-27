package providers

import "github.com/reliastra/cloudvitals/internal/core"

type Provider interface {
	Name() string
	RunCheck(check core.CheckConfig, profile string, region string) (*core.CheckResult, error)
}

