/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/util"
)

func TestFindBooleanIssues(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}

	fset := token.NewFileSet()
	var astFiles []*ast.File
	var nonTestPaths []string

	for _, fpath := range files {
		if strings.HasSuffix(fpath, "_test.go") {
			continue
		}
		src, err := os.ReadFile(fpath)
		if err != nil {
			t.Fatal(err)
		}
		node, err := parser.ParseFile(fset, fpath, src, parser.ParseComments)
		if err != nil {
			t.Fatalf("Parse error in %s: %v", fpath, err)
		}
		astFiles = append(astFiles, node)
		nonTestPaths = append(nonTestPaths, fpath)
	}

	conf := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			// ignore type checking errors if any during partial parse
		},
	}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	pkg, err := conf.Check("jacobin/src/gfunction/javaIo", fset, astFiles, info)
	if err != nil {
		t.Logf("Type check warning: %v", err)
	}
	_ = pkg

	for _, node := range astFiles {
		fpath := fset.Position(node.Pos()).Filename
		ast.Inspect(node, func(n ast.Node) bool {
			// Check for params[x].(bool) or any .(bool)
			if typeAssert, ok := n.(*ast.TypeAssertExpr); ok {
				if ident, isIdent := typeAssert.Type.(*ast.Ident); isIdent && ident.Name == "bool" {
					pos := fset.Position(typeAssert.Pos())
					t.Errorf("%s: Type assertion .(bool) found in %s", pos, fpath)
				}
			}

			// Check for functions returning any / interface{}
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Type.Results != nil && len(fn.Type.Results.List) == 1 {
					retType := fn.Type.Results.List[0].Type
					isAny := false
					if ident, isIdent := retType.(*ast.Ident); isIdent && (ident.Name == "any" || ident.Name == "interface{}") {
						isAny = true
					}
					if isAny && fn.Body != nil {
						ast.Inspect(fn.Body, func(inner ast.Node) bool {
							if nestedFn, isNested := inner.(*ast.FuncLit); isNested {
								_ = nestedFn
								return false
							}
							if ret, isRet := inner.(*ast.ReturnStmt); isRet {
								for _, res := range ret.Results {
									if isBoolExpr(res) {
										pos := fset.Position(ret.Pos())
										t.Errorf("%s: Return of Go boolean expression in gfunction %s (%s)", pos, fn.Name.Name, fpath)
									} else if tv, ok := info.Types[res]; ok {
										if basic, isBasic := tv.Type.Underlying().(*types.Basic); isBasic && basic.Info()&types.IsBoolean != 0 {
											pos := fset.Position(ret.Pos())
											t.Errorf("%s: Return of Go boolean typed expression (%s) in gfunction %s (%s)", pos, tv.Type.String(), fn.Name.Name, fpath)
										}
									}
								}
							}
							return true
						})
					}
				}
			}
			return true
		})
	}
}

func isBoolExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		if e.Name == "true" || e.Name == "false" {
			return true
		}
	case *ast.BinaryExpr:
		switch e.Op {
		case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.LAND, token.LOR:
			return true
		}
	case *ast.UnaryExpr:
		if e.Op == token.NOT {
			return true
		}
	}
	return false
}

func TestVerifyBooleanSignatures(t *testing.T) {
	loadAllGFunctions()

	countZRet := 0
	countZParam := 0
	for k := range ghelpers.MethodSignatures {
		if !strings.HasPrefix(k, "java/io/") {
			continue
		}
		_, _, methType := util.ParseFQN(k)
		if strings.HasSuffix(methType, ")Z") {
			countZRet++
		}
		params := util.ParseIncomingParamsFromMethTypeString(methType)
		for _, p := range params {
			if p == "Z" {
				countZParam++
				t.Logf("Method with Z param: %s", k)
			}
		}
	}
	t.Logf("Total java/io methods returning Z: %d, methods taking Z: %d", countZRet, countZParam)
}
