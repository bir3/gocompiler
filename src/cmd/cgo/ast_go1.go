// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build compiler_bootstrap

package cgo

import (
	"github.com/bir3/gocompiler/src/go/ast"
	"github.com/bir3/gocompiler/src/go/token"
)

func (f *File) walkUnexpected(x interface{}, context astContext, visit func(*File, interface{}, astContext)) {
	error_(token.NoPos, "unexpected type %T in walk", x)
	panic("unexpected type")
}

func funcTypeTypeParams(n *ast.FuncType) *ast.FieldList {
	return nil
}

func typeSpecTypeParams(n *ast.TypeSpec) *ast.FieldList {
	return nil
}
