// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Package parse implements the parser for the Calc programming language
package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/scan"
	"github.com/rthornton128/calc/token"
)

// ParseExpression parses the given source string and returns an ast.Node
// representing the root of the expression. This function is intended to
// facilitate testing and is not use by the compiler itself. The name is
// used in error reporting
func ParseExpression(name, src string) (ast.Expr, error) {
	var p parser

	fset := token.NewFileSet()
	file := fset.Add(name, src)
	p.init(file, name, string(src), nil)
	node := p.parseGenExpr()

	if p.errors.Count() > 0 {
		return nil, p.errors
	}
	return node, nil
}

// ParseFile parses the file identified by filename and returns a pointer
// to an ast.File object. The file should contain Calc source code and
// have the .calc file extension.
// The returned AST object ast.File is nil if there is an error.
func ParseFile(fset *token.FileSet, filename, src string) (*ast.File, error) {
	if src == "" {
		fi, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}

		if ext := filepath.Ext(fi.Name()); ext != ".calc" {
			return nil, fmt.Errorf("unknown file extension, must be .calc")
		}
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		src = string(b)
	}
	file := fset.Add(filepath.Base(filename), string(src))
	var p parser
	p.init(file, filename, string(src), ast.NewScope(nil))
	f := p.parseFile()

	if p.errors.Count() > 0 {
		return nil, p.errors
	}

	return f, nil
}

// ParseDir parses a directory of Calc source files. It calls ParseFile
// for each file ending in .calc found in the directory.
func ParseDir(fset *token.FileSet, path string) (*ast.Package, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fnames, err := fd.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	fnames = filterByExt(fnames)
	if len(fnames) == 0 {
		return nil, fmt.Errorf("no files to parse; stop")
	}

	var files []*ast.File

	for _, name := range fnames {
		f, err := ParseFile(fset, filepath.Join(path, name), "")
		if f == nil {
			return nil, err
		}
		files = append(files, f)
	}
	return &ast.Package{Files: files}, nil
}

func filterByExt(names []string) []string {
	filtered := make([]string, 0, len(names))
	for _, name := range names {
		if filepath.Ext(name) == ".calc" {
			filtered = append(filtered, name)
		}
	}
	return filtered
}

type parser struct {
	file    *token.File
	errors  token.ErrorList
	scanner scan.Scanner
	listok  bool

	curScope *ast.Scope
	topScope *ast.Scope

	pos token.Pos
	tok token.Token
	lit string
}

/* Utility */

func (p *parser) addError(args ...interface{}) {
	p.errors.Add(p.file.Position(p.pos), args...)
}

func (p *parser) checkExpr(e ast.Expr) ast.Expr {
	if e != nil && !reflect.ValueOf(e).IsNil() {
		switch t := e.(type) {
		case *ast.BasicLit, *ast.BinaryExpr, *ast.CallExpr, *ast.Ident, *ast.IfExpr,
			*ast.UnaryExpr:
		case *ast.ExprList:
			p.checkExpr(t.List[len(t.List)-1])
		default:
			p.errors.Add(p.file.Position(e.Pos()), "expression has no side-effects")
		}
	}
	return e
}

func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.addError("Expected '" + tok.String() + "' got '" + p.lit + "'")
	}
	p.next()
	return pos
}

func (p *parser) init(file *token.File, fname, src string, s *ast.Scope) {
	if s == nil {
		s = ast.NewScope(nil)
	}
	p.file = file
	p.scanner.Init(p.file, src)
	p.listok = false
	p.curScope = s
	p.topScope = s
	p.next()
}

func (p *parser) next() {
	p.lit, p.tok, p.pos = p.scanner.Scan()
}

/* Scope */

func (p *parser) openScope() {
	p.curScope = ast.NewScope(p.curScope)
}

func (p *parser) closeScope() {
	p.curScope = p.curScope.Parent
}

/* Parsing */

func (p *parser) parseAssignExpr(open token.Pos) *ast.AssignExpr {
	pos := p.expect(token.ASSIGN)
	nam := p.parseIdent()
	val := p.parseGenExpr()
	end := p.expect(token.RPAREN)

	return &ast.AssignExpr{
		Expression: ast.Expression{Opening: open, Closing: end},
		Equal:      pos,
		Name:       nam,
		Value:      p.checkExpr(val),
	}
}

func (p *parser) parseBasicLit() *ast.BasicLit {
	pos, tok, lit := p.pos, p.tok, p.lit
	p.next()
	return &ast.BasicLit{LitPos: pos, Kind: tok, Lit: lit}
}

func (p *parser) parseBinaryExpr(open token.Pos) *ast.BinaryExpr {
	pos := p.pos
	op := p.tok
	p.next()

	var list []ast.Expr
	for p.tok != token.RPAREN && p.tok != token.EOF {
		list = append(list, p.parseGenExpr())
	}
	if len(list) < 2 {
		p.addError("binary expression must have at least two operands")
	}
	end := p.expect(token.RPAREN)
	return &ast.BinaryExpr{
		Expression: ast.Expression{
			Opening: open,
			Closing: end,
		},
		Op:    op,
		OpPos: pos,
		List:  list,
	}
}

func (p *parser) parseCallExpr(open token.Pos) *ast.CallExpr {
	nam := p.parseIdent()

	var list []ast.Expr
	for p.tok != token.RPAREN && p.tok != token.EOF {
		list = append(list, p.parseGenExpr())
	}
	end := p.expect(token.RPAREN)

	return &ast.CallExpr{
		Expression: ast.Expression{
			Opening: open,
			Closing: end,
		},
		Name: nam,
		Args: list,
	}
}

func (p *parser) parseDeclExpr(open token.Pos) *ast.DeclExpr {
	if p.curScope != p.topScope {
		p.addError("function declarations may only be used in top-level scope")
		return nil
	}
	pos := p.expect(token.DECL)
	name := p.parseIdent()

	p.openScope()

	var list []*ast.Ident
	if p.tok == token.LPAREN {
		p.next()
		list = p.parseParamList()
	}

	typ := p.parseIdent()

	body := p.tryExprOrList()
	end := p.expect(token.RPAREN)

	decl := &ast.DeclExpr{
		Expression: ast.Expression{
			Opening: open,
			Closing: end,
		},
		Decl:   pos,
		Name:   name,
		Type:   typ,
		Params: list,
		Body:   p.checkExpr(body),
	}

	p.closeScope()

	prev := p.curScope.Insert(&ast.Object{
		NamePos: name.NamePos,
		Name:    name.Name,
		Kind:    ast.FuncDecl,
	})
	if prev != nil {
		switch prev.Kind {
		case ast.FuncDecl:
			p.addError(name.Name, " redeclared; declared as function at ",
				p.file.Position(prev.NamePos))
		case ast.VarDecl:
			p.addError(name.Name, " redeclared; declared as variable at ",
				p.file.Position(prev.NamePos))
		}
	}

	return decl
}

func (p *parser) parseExpr() ast.Expr {
	var expr ast.Expr
	listok := p.listok

	pos := p.expect(token.LPAREN)
	if p.listok && p.tok == token.LPAREN {
		expr = p.parseExprList(pos)
		return expr
	}
	p.listok = false
	switch p.tok {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM,
		token.EQL, token.GTE, token.GTT, token.NEQ, token.LST, token.LTE:
		expr = p.parseBinaryExpr(pos)
	case token.ASSIGN:
		expr = p.parseAssignExpr(pos)
	case token.DECL:
		expr = p.parseDeclExpr(pos)
	case token.IDENT:
		expr = p.parseCallExpr(pos)
	case token.IF:
		expr = p.parseIfExpr(pos)
	case token.VAR:
		expr = p.parseVarExpr(pos)
	default:
		if listok {
			p.addError("Expected expression but got '" + p.lit + "'")
		} else {
			p.addError("Expected operator, keyword or identifier but got '" + p.lit +
				"'")
		}
	}

	return expr
}

func (p *parser) parseExprList(open token.Pos) ast.Expr {
	p.listok = false
	var list []ast.Expr
	for p.tok != token.RPAREN {
		list = append(list, p.parseGenExpr())
	}
	if len(list) < 1 {
		p.addError("empty expression list not allowed")
	}
	end := p.expect(token.RPAREN)
	return &ast.ExprList{
		Expression: ast.Expression{
			Opening: open,
			Closing: end,
		},
		List: list,
	}
}

func (p *parser) parseGenExpr() ast.Expr {
	var expr ast.Expr

	switch p.tok {
	case token.LPAREN:
		expr = p.parseExpr()
	case token.IDENT:
		expr = p.parseIdent()
	case token.BOOL, token.INTEGER:
		expr = p.parseBasicLit()
	case token.ADD, token.SUB:
		expr = p.parseUnaryExpr()
	default:
		p.addError("Expected expression, got '" + p.lit + "'")
		p.next()
	}
	p.listok = false

	return expr
}

func (p *parser) parseFile() *ast.File {
	decls := make([]*ast.DeclExpr, 0, 16)
	for p.tok != token.EOF {
		e := p.parseGenExpr()
		if decl, ok := e.(*ast.DeclExpr); ok {
			decls = append(decls, decl)
		}
	}
	if len(decls) < 1 {
		p.addError("reached end of file without any declarations")
	}
	return &ast.File{Decls: decls}
}

func (p *parser) parseIdent() *ast.Ident {
	name := p.lit
	return &ast.Ident{NamePos: p.expect(token.IDENT), Name: name}
}

func (p *parser) parseIfExpr(open token.Pos) *ast.IfExpr {
	pos := p.expect(token.IF)
	cond := p.parseGenExpr()
	typ := p.parseIdent()

	p.openScope()
	then := p.tryExprOrList()
	var els ast.Expr
	if p.tok != token.RPAREN {
		els = p.tryExprOrList()
	}
	p.closeScope()
	end := p.expect(token.RPAREN)

	return &ast.IfExpr{
		Expression: ast.Expression{
			Opening: open,
			Closing: end,
		},
		If:   pos,
		Type: typ,
		Cond: cond,
		Then: then,
		Else: els,
	}
}

func (p *parser) parseParamList() []*ast.Ident {
	var list []*ast.Ident
	count, start := 0, 0
	for p.tok != token.RPAREN {
		ident := p.parseIdent()
		count++
		if p.tok == token.COMMA || p.tok == token.RPAREN {
			for _, param := range list[start:] {
				o := &ast.Object{
					Kind:    ast.VarDecl,
					Name:    param.Name,
					NamePos: param.Pos(),
				}
				param.Type = ident
				if prev := p.curScope.Insert(o); prev != nil {
					p.addError("duplicate parameter ", param.Name,
						"; previously declared at ", p.file.Position(prev.Pos()))
				}
			}
			start = count
			continue
		}
		list = append(list, ident)
	}
	if len(list) < 1 {
		p.addError("empty param list not allowed")
	}
	p.expect(token.RPAREN)
	return list
}

func (p *parser) parseUnaryExpr() *ast.UnaryExpr {
	pos, op := p.pos, p.lit
	p.next()
	exp := p.parseGenExpr()
	return &ast.UnaryExpr{OpPos: pos, Op: op, Value: p.checkExpr(exp)}
}

func (p *parser) parseVarExpr(open token.Pos) *ast.VarExpr {
	varpos := p.expect(token.VAR)

	var name *ast.Ident
	var value *ast.AssignExpr
	switch p.tok {
	case token.IDENT:
		name = p.parseIdent()
	case token.LPAREN:
		value = p.parseAssignExpr(p.expect(token.LPAREN))
		name = value.Name
	default:
		p.addError("expected identifier or assignment")
		name = &ast.Ident{NamePos: token.NoPos, Name: "NoName"}
	}
	if value == nil || p.tok == token.IDENT {
		name.Type = p.parseIdent()
	}
	end := p.expect(token.RPAREN)

	prev := p.curScope.Insert(&ast.Object{
		NamePos: name.NamePos,
		Name:    name.Name,
		Kind:    ast.VarDecl,
	})
	if prev != nil {
		p.addError(name.Name, " redeclared; declared as ", prev.Kind.String(),
			" at ", p.file.Position(prev.NamePos))

	}

	return &ast.VarExpr{
		Expression: ast.Expression{Opening: open, Closing: end},
		Var:        varpos,
		Name:       name,
		Value:      value,
	}
}

func (p *parser) tryExprOrList() ast.Expr {
	p.listok = true
	return p.parseGenExpr()
}
