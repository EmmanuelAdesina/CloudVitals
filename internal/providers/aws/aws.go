package aws

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/reliastra/cloudvitals/internal/core"
)

type AWSProvider struct {
	scriptBase string
}

func New() *AWSProvider {
	return &AWSProvider{
		scriptBase: filepath.Join("internal", "providers", "aws", "checks"),
	}
}

func (a *AWSProvider) Name() string {
	return "aws"
}

func (a *AWSProvider) RunCheck(check core.CheckConfig, profile string, region string) (*core.CheckResult, error) {
	scriptPath := filepath.Join(a.scriptBase, check.Script)

	cmd := exec.Command("python3", scriptPath,
		"--profile", profile,
		"--region", region,
	)

	output, err := cmd.CombinedOutput()

	var result core.CheckResult
	if parseErr := json.Unmarshal(output, &result); parseErr != nil {
		return nil, fmt.Errorf("check %s failed: %v | output: %s", check.ID, err, string(output))
	}

	// Enrich with registry metadata
	result.CheckID = check.ID
	result.CheckName = check.Name
	result.Severity = check.Severity
	result.ExecutedAt = time.Now().UTC()

	return &result, nil
}
