package main

import (
	"github.com/romariok/golanglog-linter/pkg/golanglog"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(golanglog.Analyzer)
}
