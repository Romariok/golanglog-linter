package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/romariok/golanglog-linter/pkg/golanglog/config"
)

type keywordPattern struct {
	name string
	re   *regexp.Regexp
}

// CheckSensitive reports a diagnostic when a log call may expose sensitive data:
//   - a string literal in any argument contains a sensitive keyword (word-boundary match)
//   - an identifier or selector field name in any argument matches a sensitive keyword
func CheckSensitive(pass *analysis.Pass, call *ast.CallExpr, msgArg ast.Expr, cfg *config.Config) {
	patterns := buildSensitivePatterns(cfg)
	if len(patterns) == 0 {
		return
	}

	reported := false
	for _, arg := range call.Args {
		walkSensitive(pass, arg, patterns, &reported)
		if reported {
			return
		}
	}
}

func buildSensitivePatterns(cfg *config.Config) []keywordPattern {
	var patterns []keywordPattern
	for _, kw := range cfg.SensitiveKeywords {
		re := regexp.MustCompile(fmt.Sprintf(`(?i)(^|[^a-zA-Z])%s([^a-zA-Z]|$)`, regexp.QuoteMeta(kw)))
		patterns = append(patterns, keywordPattern{name: kw, re: re})
	}
	for _, pat := range cfg.CustomPatterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			continue
		}
		patterns = append(patterns, keywordPattern{name: pat, re: re})
	}
	return patterns
}

func walkSensitive(pass *analysis.Pass, expr ast.Expr, patterns []keywordPattern, reported *bool) {
	if *reported {
		return
	}
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return
		}
		s, err := strconv.Unquote(e.Value)
		if err != nil {
			return
		}
		for _, p := range patterns {
			if p.re.MatchString(s) {
				pass.Reportf(e.Pos(), "log message may contain sensitive data: found keyword %q", p.name)
				*reported = true
				return
			}
		}

	case *ast.Ident:
		lower := strings.ToLower(e.Name)
		for _, p := range patterns {
			if strings.Contains(lower, strings.ToLower(p.name)) {
				pass.Reportf(e.Pos(), "log message may contain sensitive data: found keyword %q", p.name)
				*reported = true
				return
			}
		}

	case *ast.SelectorExpr:
		walkSensitive(pass, e.X, patterns, reported)
		if !*reported {
			lower := strings.ToLower(e.Sel.Name)
			for _, p := range patterns {
				if strings.Contains(lower, strings.ToLower(p.name)) {
					pass.Reportf(e.Sel.Pos(), "log message may contain sensitive data: found keyword %q", p.name)
					*reported = true
					return
				}
			}
		}

	case *ast.BinaryExpr:
		walkSensitive(pass, e.X, patterns, reported)
		if !*reported {
			walkSensitive(pass, e.Y, patterns, reported)
		}

	case *ast.CallExpr:
		for _, arg := range e.Args {
			walkSensitive(pass, arg, patterns, reported)
			if *reported {
				return
			}
		}

	case *ast.ParenExpr:
		walkSensitive(pass, e.X, patterns, reported)
	}
}
