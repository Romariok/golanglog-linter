//go:build ignore

// Plugin is the golangci-lint module plugin entry point.
// Build with: go build -buildmode=plugin -o golanglog.so ./plugin/
package main

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
