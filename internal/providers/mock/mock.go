package mock

import (
	"time"

	"github.com/reliastra/cloudvitals/internal/core"
)

type MockProvider struct{}

func New() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) Name() string {
	return "mock"
}

func (m *MockProvider) RunCheck(check core.CheckConfig, profile string, region string) (*core.CheckResult, error) {
	// Simulate: s3_public fails, everything else passes
	status := "pass"
	findings := []core.Finding{}

	if check.ID == "s3_public" {
		status = "fail"
		findings = append(findings, core.Finding{
			Resource:    "arn:aws:s3:::test-bucket",
			Region:      region,
			Detail:      "S3 bucket public access block not fully enabled",
			Remediation: "aws s3api put-public-access-block --bucket test-bucket ...",
		})
	}

	return &core.CheckResult{
		CheckID:    check.ID,
		CheckName:  check.Name,
		Status:     status,
		Severity:   check.Severity,
		Findings:   findings,
		ExecutedAt: time.Now().UTC(),
	}, nil
}
