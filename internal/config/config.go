package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RuleConfig struct {
	Level  string         `yaml:"level,omitempty"`
	Config map[string]any `yaml:",inline"`
}

type Config struct {
	Extends string                `yaml:"extends,omitempty"`
	Rules   map[string]RuleConfig `yaml:"rules"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		Rules: map[string]RuleConfig{
			"anchors":              {Level: "enable"},
			"braces":               {Level: "enable"},
			"brackets":             {Level: "enable"},
			"colons":               {Level: "enable"},
			"commas":               {Level: "enable"},
			"comments":             {Level: "enable"},
			"comments-indentation": {Level: "enable"},
			"document-end":         {Level: "enable"},
			"document-start":       {Level: "enable"},
			"empty-lines":          {Level: "enable"},
			"empty-values":         {Level: "enable"},
			"float-values":         {Level: "enable"},
			"hyphens":              {Level: "enable"},
			"indentation":          {Level: "enable"},
			"key-duplicates":       {Level: "enable"},
			"key-ordering":         {Level: "disable"},
			"line-length":          {Level: "enable"},
			"new-lines":            {Level: "enable"},
			"octal-values":         {Level: "enable"},
			"quoted-strings":       {Level: "disable"},
			"trailing-lines":       {Level: "enable"},
			"trailing-spaces":      {Level: "enable"},
			"truthy":               {Level: "enable"},
		},
	}
}
