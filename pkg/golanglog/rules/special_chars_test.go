package rules

import (
	"go/ast"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var specialCharsAnalyzer = &analysis.Analyzer{
	Name:     "specialchars",
	Doc:      "check that log messages do not contain special characters or emojis",
	Run:      runSpecialCharsOnly,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runSpecialCharsOnly(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		msgArg, ok := IsLogCall(pass, call)
		if !ok {
			return
		}
		CheckSpecialChars(pass, call, msgArg)
	})
	return nil, nil
}

func TestSpecialChars(t *testing.T) {
	analysistest.Run(t, testdataDir(), specialCharsAnalyzer, "special")
}
