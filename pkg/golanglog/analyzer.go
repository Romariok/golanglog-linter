package golanglog

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/romariok/golanglog-linter/pkg/golanglog/config"
	"github.com/romariok/golanglog-linter/pkg/golanglog/rules"
)

const doc = `golanglog validates log messages for style and security issues.

Supported loggers: log/slog, go.uber.org/zap (Logger and SugaredLogger).

Rules:
  lowercase     — message must start with a lowercase letter (provides SuggestedFix)
  english       — message must contain only ASCII characters
  special-chars — message must not contain emojis, repeated/trailing special chars, or \n
  sensitive     — message must not reference sensitive data keywords
`

var cfg = config.Default()

var Analyzer = &analysis.Analyzer{
	Name:     "golanglog",
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() {
	Analyzer.Flags.BoolVar(&cfg.Rules.Lowercase, "lowercase", cfg.Rules.Lowercase,
		"check that log messages start with a lowercase letter")
	Analyzer.Flags.BoolVar(&cfg.Rules.English, "english", cfg.Rules.English,
		"check that log messages contain only ASCII characters")
	Analyzer.Flags.BoolVar(&cfg.Rules.SpecialChars, "special-chars", cfg.Rules.SpecialChars,
		"check that log messages do not contain special characters or emojis")
	Analyzer.Flags.BoolVar(&cfg.Rules.Sensitive, "sensitive", cfg.Rules.Sensitive,
		"check that log messages do not contain sensitive data keywords")
	Analyzer.Flags.Var(
		&config.StringSliceFlag{Values: &cfg.SensitiveKeywords},
		"sensitive-keywords",
		"comma-separated list of sensitive keywords to detect",
	)
	Analyzer.Flags.Var(
		&config.StringSliceFlag{Values: &cfg.CustomPatterns},
		"custom-patterns",
		"comma-separated list of custom regex patterns for sensitive data detection",
	)
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		msgArg, ok := rules.IsLogCall(pass, call)
		if !ok {
			return
		}

		if cfg.Rules.Lowercase {
			rules.CheckLowercase(pass, call, msgArg)
		}
		if cfg.Rules.English {
			rules.CheckEnglish(pass, call, msgArg)
		}
		if cfg.Rules.SpecialChars {
			rules.CheckSpecialChars(pass, call, msgArg)
		}
		if cfg.Rules.Sensitive {
			rules.CheckSensitive(pass, call, msgArg, cfg)
		}
	})

	return nil, nil
}
