package config

import "strings"

// Config holds all configuration for the golanglog-linter analyzer.
type Config struct {
	Rules             RulesConfig
	SensitiveKeywords []string
	CustomPatterns    []string
}

// RulesConfig controls which rules are enabled.
type RulesConfig struct {
	Lowercase    bool
	English      bool
	SpecialChars bool
	Sensitive    bool
}

// Default returns a Config with all rules enabled and the default sensitive keywords list.
func Default() *Config {
	return &Config{
		Rules: RulesConfig{
			Lowercase:    true,
			English:      true,
			SpecialChars: true,
			Sensitive:    true,
		},
		SensitiveKeywords: []string{
			"password", "passwd", "secret", "token",
			"api_key", "apikey", "auth", "credential", "private_key",
		},
	}
}

// StringSliceFlag implements flag.Value for a comma-separated list of strings.
type StringSliceFlag struct {
	Values *[]string
}

func (f *StringSliceFlag) String() string {
	if f.Values == nil || *f.Values == nil {
		return ""
	}
	return strings.Join(*f.Values, ",")
}

func (f *StringSliceFlag) Set(s string) error {
	if s == "" {
		*f.Values = nil
		return nil
	}
	*f.Values = strings.Split(s, ",")
	return nil
}
