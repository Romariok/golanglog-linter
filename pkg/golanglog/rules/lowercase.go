package rules

import (
	"go/ast"
	"go/token"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

// CheckLowercase reports a diagnostic when the log message string literal starts
// with an uppercase Unicode letter. It also provides a SuggestedFix.
func CheckLowercase(pass *analysis.Pass, call *ast.CallExpr, msgArg ast.Expr) {
	str, lit, ok := extractStringLit(msgArg)
	if !ok || str == "" {
		return
	}

	firstRune, runeLen := utf8.DecodeRuneInString(str)
	if firstRune == utf8.RuneError || !unicode.IsUpper(firstRune) {
		return
	}

	// firstCharPos is the source position of the first character (after the opening quote).
	firstCharPos := lit.Pos() + 1

	pass.Report(analysis.Diagnostic{
		Pos:     lit.Pos(),
		Message: "log message should start with a lowercase letter",
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: "Convert first letter to lowercase",
			TextEdits: []analysis.TextEdit{{
				Pos:     firstCharPos,
				End:     firstCharPos + token.Pos(runeLen),
				NewText: []byte(string(unicode.ToLower(firstRune))),
			}},
		}},
	})
}
