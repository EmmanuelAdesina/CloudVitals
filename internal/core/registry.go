package core

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Registry struct {
	Checks []CheckConfig `yaml:"checks"`
}

func LoadRegistry(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var r Registry
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *Registry) ChecksForProvider(provider string) []CheckConfig {
	var filtered []CheckConfig
	target := strings.ToLower(provider)
	for _, c := range r.Checks {
		if strings.ToLower(c.Provider) == target {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
