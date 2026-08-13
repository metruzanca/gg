package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Types []CommitType `yaml:"types"`
}

type CommitType struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

var DefaultTypes = []CommitType{
	{Name: "fix", Description: "bug fix"},
	{Name: "feat", Description: "new feature"},
	{Name: "chore", Description: "maintenance"},
	{Name: "docs", Description: "documentation"},
	{Name: "refactor", Description: "code restructuring"},
	{Name: "test", Description: "adding tests"},
	{Name: "style", Description: "formatting"},
	{Name: "perf", Description: "performance improvement"},
	{Name: "ci", Description: "CI/CD changes"},
	{Name: "build", Description: "build system changes"},
}

func Load(path string) (*Config, error) {
	cfg := &Config{Types: DefaultTypes}
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return cfg, nil
		}
		path = findConfig(cwd)
	}
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil
	}
	var userCfg Config
	if err := yaml.Unmarshal(data, &userCfg); err != nil {
		return cfg, nil
	}
	if len(userCfg.Types) > 0 {
		cfg.Types = userCfg.Types
	}
	return cfg, nil
}

func findConfig(dir string) string {
	for {
		p := filepath.Join(dir, ".gg.yaml")
		if _, err := os.Stat(p); err == nil {
			return p
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
