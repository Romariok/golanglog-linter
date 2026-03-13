package rules

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const (
	slogPkgPath    = "log/slog"
	zapPkgPath     = "go.uber.org/zap"
	zapLoggerName  = "Logger"
	zapSugaredName = "SugaredLogger"
)

// slogMethodArgs maps slog function/method names to the 0-based index of the message argument.
var slogMethodArgs = map[string]int{
	"Debug":        0,
	"Info":         0,
	"Warn":         0,
	"Error":        0,
	"DebugContext": 1,
	"InfoContext":  1,
	"WarnContext":  1,
	"ErrorContext": 1,
	"Log":          2, // Log(ctx, level, msg, ...)
	"LogAttrs":     2, // LogAttrs(ctx, level, msg, ...)
}

// zapLoggerArgs maps (*zap.Logger) method names to the 0-based message argument index.
var zapLoggerArgs = map[string]int{
	"Debug":  0,
	"Info":   0,
	"Warn":   0,
	"Error":  0,
	"DPanic": 0,
	"Panic":  0,
	"Fatal":  0,
}

// zapSugarArgs maps (*zap.SugaredLogger) method names to the 0-based message argument index.
var zapSugarArgs = map[string]int{
	"Debug":   0,
	"Info":    0,
	"Warn":    0,
	"Error":   0,
	"DPanic":  0,
	"Panic":   0,
	"Fatal":   0,
	"Debugf":  0,
	"Infof":   0,
	"Warnf":   0,
	"Errorf":  0,
	"DPanicf": 0,
	"Panicf":  0,
	"Fatalf":  0,
}

// IsLogCall reports whether call is an invocation of a supported logging function.
// If it is, it returns the message argument expression and true.
// Detection is done via pass.TypesInfo to check the exact package/type of the callee,
// avoiding false positives from methods with the same name in other packages.
func IsLogCall(pass *analysis.Pass, call *ast.CallExpr) (ast.Expr, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	methodName := sel.Sel.Name

	obj, ok := pass.TypesInfo.Uses[sel.Sel]
	if !ok || obj == nil {
		return nil, false
	}

	fn, ok := obj.(*types.Func)
	if !ok || fn.Pkg() == nil {
		return nil, false
	}

	pkgPath := fn.Pkg().Path()

	switch pkgPath {
	case slogPkgPath:
		// Covers both package-level slog.Info(...) and (*slog.Logger).Info(...).
		if idx, ok := slogMethodArgs[methodName]; ok && len(call.Args) > idx {
			return call.Args[idx], true
		}

	case zapPkgPath:
		sig, ok := fn.Type().(*types.Signature)
		if !ok {
			return nil, false
		}
		recv := sig.Recv()
		if recv == nil {
			return nil, false
		}
		// Dereference pointer receiver (*Logger → Logger, *SugaredLogger → SugaredLogger).
		recvType := recv.Type()
		if ptr, ok := recvType.(*types.Pointer); ok {
			recvType = ptr.Elem()
		}
		named, ok := recvType.(*types.Named)
		if !ok {
			return nil, false
		}
		switch named.Obj().Name() {
		case zapLoggerName:
			if idx, ok := zapLoggerArgs[methodName]; ok && len(call.Args) > idx {
				return call.Args[idx], true
			}
		case zapSugaredName:
			if idx, ok := zapSugarArgs[methodName]; ok && len(call.Args) > idx {
				return call.Args[idx], true
			}
		}
	}

	return nil, false
}
