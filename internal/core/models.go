package core

import "time"

type Severity string

const (
	Critical Severity = "critical"
	High     Severity = "high"
	Medium   Severity = "medium"
	Low      Severity = "low"
)

type CheckConfig struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Provider    string   `yaml:"provider"`
	Severity    Severity `yaml:"severity"`
	Script      string   `yaml:"script"`
}

type Finding struct {
	Resource    string `json:"resource"`
	Region      string `json:"region"`
	Detail      string `json:"detail"`
	Remediation string `json:"remediation"`
}

type CheckResult struct {
	CheckID    string    `json:"check_id"`
	CheckName  string    `json:"check_name"`
	Status     string    `json:"status"`
	Severity   Severity  `json:"severity"`
	Findings   []Finding `json:"findings"`
	ExecutedAt time.Time `json:"executed_at"`
	ErrorMsg   string    `json:"error_msg,omitempty"`
}

type Provider interface {
	Name() string
	RunCheck(check CheckConfig, profile string, region string) (*CheckResult, error)
}
