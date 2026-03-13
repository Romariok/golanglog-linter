// Package plugin is the golangci-lint module plugin entry point.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/romariok/golanglog-linter/pkg/golanglog"
	"github.com/romariok/golanglog-linter/pkg/golanglog/config"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("golanglog", New)
}

type pluginSettings struct {
	Rules struct {
		Lowercase    bool `json:"lowercase"`
		English      bool `json:"english"`
		SpecialChars bool `json:"special-chars"`
		Sensitive    bool `json:"sensitive"`
	} `json:"rules"`
	SensitiveKeywords []string `json:"sensitive-keywords"`
	CustomPatterns    []string `json:"custom-patterns"`
}

type linterPlugin struct {
	cfg *config.Config
}

// New creates a new plugin instance. conf is nil when no settings block is provided.
func New(conf any) (register.LinterPlugin, error) {
	cfg := config.Default()
	if conf != nil {
		s, err := register.DecodeSettings[pluginSettings](conf)
		if err != nil {
			return nil, err
		}
		cfg.Rules.Lowercase = s.Rules.Lowercase
		cfg.Rules.English = s.Rules.English
		cfg.Rules.SpecialChars = s.Rules.SpecialChars
		cfg.Rules.Sensitive = s.Rules.Sensitive
		if len(s.SensitiveKeywords) > 0 {
			cfg.SensitiveKeywords = s.SensitiveKeywords
		}
		if len(s.CustomPatterns) > 0 {
			cfg.CustomPatterns = s.CustomPatterns
		}
	}
	return &linterPlugin{cfg: cfg}, nil
}

func (p *linterPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{golanglog.NewAnalyzer(p.cfg)}, nil
}

func (p *linterPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
