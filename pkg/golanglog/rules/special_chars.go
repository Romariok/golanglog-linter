package rules

import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var (
	specialCharsRe = regexp.MustCompile(`[.!?]{2,}|[?:]$|!`)
	emojiRe = regexp.MustCompile(`[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}]`)
)

// CheckSpecialChars reports a diagnostic when the log message string literal
// contains emojis, repeated special characters, a trailing !/?/:, or \n.
func CheckSpecialChars(pass *analysis.Pass, call *ast.CallExpr, msgArg ast.Expr) {
	reported := false
	forEachStringLit(msgArg, func(val string, lit *ast.BasicLit) {
		if reported {
			return
		}
		if strings.ContainsRune(val, '\n') || emojiRe.MatchString(val) || specialCharsRe.MatchString(val) {
			pass.Reportf(lit.Pos(), "log message must not contain special characters or emojis")
			reported = true
		}
	})
}
