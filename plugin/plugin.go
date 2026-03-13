// Package plugin is the golangci-lint module plugin entry point.
package plugin

import (
	"github.com/romariok/golanglog-linter/pkg/golanglog"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{golanglog.Analyzer}
}

// AnalyzerPlugin is the exported plugin symbol for golangci-lint module plugin.
var AnalyzerPlugin analyzerPlugin
