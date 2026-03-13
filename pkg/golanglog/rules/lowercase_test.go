package rules

import (
	"go/ast"
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "..", "testdata")
}

var lowercaseAnalyzer = &analysis.Analyzer{
	Name:     "lowercase",
	Doc:      "check that log messages start with a lowercase letter",
	Run:      runLowercaseOnly,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runLowercaseOnly(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		msgArg, ok := IsLogCall(pass, call)
		if !ok {
			return
		}
		CheckLowercase(pass, call, msgArg)
	})
	return nil, nil
}

func TestLowercase(t *testing.T) {
	analysistest.Run(t, testdataDir(), lowercaseAnalyzer, "lowercase")
}
