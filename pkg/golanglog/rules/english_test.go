package rules

import (
	"go/ast"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var englishAnalyzer = &analysis.Analyzer{
	Name:     "english",
	Doc:      "check that log messages contain only ASCII characters",
	Run:      runEnglishOnly,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runEnglishOnly(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		msgArg, ok := IsLogCall(pass, call)
		if !ok {
			return
		}
		CheckEnglish(pass, call, msgArg)
	})
	return nil, nil
}

func TestEnglish(t *testing.T) {
	analysistest.Run(t, testdataDir(), englishAnalyzer, "english")
}
