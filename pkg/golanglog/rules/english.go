package rules

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// CheckEnglish reports a diagnostic when the log message string literal contains
// non-ASCII characters (rune > 127).
func CheckEnglish(pass *analysis.Pass, call *ast.CallExpr, msgArg ast.Expr) {
	reported := false
	forEachStringLit(msgArg, func(val string, lit *ast.BasicLit) {
		if reported {
			return
		}
		for _, r := range val {
			if r > 127 {
				pass.Reportf(lit.Pos(), "log message must be in English only")
				reported = true
				return
			}
		}
	})
}
