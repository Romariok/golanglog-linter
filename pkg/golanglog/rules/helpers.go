package rules

import (
	"go/ast"
	"go/token"
	"strconv"
)

// extractStringLit returns the unquoted value and AST node of the leftmost
// string literal in expr, following + chains and parentheses.
func extractStringLit(expr ast.Expr) (string, *ast.BasicLit, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			if err != nil {
				return "", nil, false
			}
			return s, e, true
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return extractStringLit(e.X)
		}
	case *ast.ParenExpr:
		return extractStringLit(e.X)
	}
	return "", nil, false
}

// forEachStringLit calls fn for every string literal reachable in expr
// via + chains and parentheses.
func forEachStringLit(expr ast.Expr, fn func(val string, lit *ast.BasicLit)) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			if err == nil {
				fn(s, e)
			}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			forEachStringLit(e.X, fn)
			forEachStringLit(e.Y, fn)
		}
	case *ast.ParenExpr:
		forEachStringLit(e.X, fn)
	}
}
