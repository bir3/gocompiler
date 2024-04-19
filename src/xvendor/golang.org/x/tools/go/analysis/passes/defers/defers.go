// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defers

import (
	_ "embed"
	"github.com/bir3/gocompiler/src/go/ast"

	"github.com/bir3/gocompiler/src/xvendor/golang.org/x/tools/go/analysis"
	"github.com/bir3/gocompiler/src/xvendor/golang.org/x/tools/go/analysis/passes/inspect"
	"github.com/bir3/gocompiler/src/xvendor/golang.org/x/tools/go/analysis/passes/internal/analysisutil"
	"github.com/bir3/gocompiler/src/xvendor/golang.org/x/tools/go/ast/inspector"
	"github.com/bir3/gocompiler/src/xvendor/golang.org/x/tools/go/types/typeutil"
)

//go:embed doc.go
var doc string

// Analyzer is the defers analyzer.
var Analyzer = &analysis.Analyzer{
	Name:		"defers",
	Requires:	[]*analysis.Analyzer{inspect.Analyzer},
	URL:		"https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/defers",
	Doc:		analysisutil.MustExtractDoc(doc, "defers"),
	Run:		run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !analysisutil.Imports(pass.Pkg, "time") {
		return nil, nil
	}

	checkDeferCall := func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.CallExpr:
			if analysisutil.IsFunctionNamed(typeutil.StaticCallee(pass.TypesInfo, v), "time", "Since") {
				pass.Reportf(v.Pos(), "call to time.Since is not deferred")
			}
		case *ast.FuncLit:
			return false	// prune
		}
		return true
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.DeferStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		d := n.(*ast.DeferStmt)
		ast.Inspect(d.Call, checkDeferCall)
	})

	return nil, nil
}
