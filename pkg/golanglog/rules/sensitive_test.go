package rules

import (
	"go/ast"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/romariok/golanglog-linter/pkg/golanglog/config"
)

var sensitiveAnalyzer = &analysis.Analyzer{
	Name:     "sensitive",
	Doc:      "check that log messages do not reference sensitive data keywords",
	Run:      runSensitiveOnly,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runSensitiveOnly(pass *analysis.Pass) (interface{}, error) {
	cfg := config.Default()
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		msgArg, ok := IsLogCall(pass, call)
		if !ok {
			return
		}
		CheckSensitive(pass, call, msgArg, cfg)
	})
	return nil, nil
}

func TestSensitive(t *testing.T) {
	analysistest.Run(t, testdataDir(), sensitiveAnalyzer, "sensitive")
}
