// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/scan"
	"github.com/rthornton128/calc/token"
)

func ParseExpression(name, src string) (node ast.Node) {
	var p parser

	fset := token.NewFileSet()
	file := fset.Add(name, src)
	p.init(file, name, string(src))
	node = p.parseGenExpr()

	if p.errors.Count() > 0 {
		p.errors.Print()
		node = nil
	}
	return
}

func ParseFile(fset *token.FileSet, filename string) *ast.File {
	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var p parser
	var f *ast.File
	if ext := filepath.Ext(fi.Name()); ext == ".calc" {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		file := fset.Add(filepath.Base(filename), string(src))
		p.init(file, filename, string(src))
		f = p.parseFile()
	}

	if p.errors.Count() > 0 {
		p.errors.Print()
		return nil
	}

	return f
}

func ParseDir(fset *token.FileSet, path string) *ast.Package {
	fd, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer fd.Close()

	fis, err := fd.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var files []*ast.File
	// TODO: use concurrency
	for _, fi := range fis {
		f := ParseFile(fset, filepath.Join(path, fi.Name()))
		if f == nil {
			return nil
		}
		files = append(files, f)
	}
	scope := ast.MergeScopes(files)
	return &ast.Package{Scope: scope, Files: files}
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
	if p.errors.Count() >= 10 {
		p.errors.Print()
		os.Exit(1)
	}
}

func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.addError("Expected '" + tok.String() + "' got '" + p.lit + "'")
	}
	p.next()
	return pos
}

func (p *parser) init(file *token.File, fname, src string) {
	p.file = file
	p.scanner.Init(p.file, src)
	p.listok = false
	p.curScope = ast.NewScope(nil)
	p.topScope = p.curScope
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
		Value:      val,
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
	nam := p.parseIdent()

	p.openScope()

	var list []*ast.Ident
	if p.tok == token.LPAREN {
		p.next()
		list = p.parseParamList()
	}

	typ := p.parseIdent()

	var body ast.Expr
	if p.tok == token.LPAREN {
		body = p.tryExprOrList()
	}
	end := p.expect(token.RPAREN)

	decl := &ast.DeclExpr{
		Expression: ast.Expression{
			Opening: open,
			Closing: end,
		},
		Decl:   pos,
		Name:   nam,
		Type:   typ,
		Params: list,
		Body:   body,
		Scope:  p.curScope,
	}
	ob := &ast.Object{
		NamePos: nam.NamePos,
		Name:    nam.Name,
		Kind:    ast.Decl,
		Type:    typ,
		Value:   decl,
	}

	p.closeScope()

	if old := p.curScope.Insert(ob); old != nil {
		p.addError("redeclaration of function not allowed, originally declared "+
			"at: ", p.file.Position(old.NamePos))
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
	case token.INTEGER:
		expr = p.parseBasicLit()
	default:
		p.addError("Expected expression, got '" + p.lit + "'")
		p.next()
	}
	p.listok = false

	return expr
}

func (p *parser) parseFile() *ast.File {
	for p.tok != token.EOF {
		p.parseGenExpr()
	}
	if p.topScope.Size() < 1 {
		p.addError("reached end of file without any declarations")
	}
	return &ast.File{Scope: p.topScope}
}

func (p *parser) parseIdent() *ast.Ident {
	name := p.lit
	pos := p.expect(token.IDENT)
	return &ast.Ident{NamePos: pos, Name: name}
}

func (p *parser) parseIfExpr(open token.Pos) *ast.IfExpr {
	pos := p.expect(token.IF)
	cond := p.parseGenExpr()

	var typ *ast.Ident
	if p.tok == token.IDENT {
		typ = p.parseIdent()
	}

	p.openScope()
	scope := p.curScope
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
		If:    pos,
		Type:  typ,
		Cond:  cond,
		Then:  then,
		Else:  els,
		Scope: scope,
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
				if param.Object == nil {
					param.Object = &ast.Object{
						Kind: ast.Var,
						Name: param.Name,
					}
				}
				param.Object.Type = ident
				p.curScope.Insert(param.Object)
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

func (p *parser) parseVarExpr(open token.Pos) *ast.VarExpr {
	var (
		name  *ast.Ident
		vtype *ast.Ident
		value *ast.AssignExpr
	)
	varpos := p.expect(token.VAR)

	switch p.tok {
	case token.IDENT:
		name = p.parseIdent()
	case token.LPAREN:
		value = p.parseAssignExpr(p.expect(token.LPAREN))
		name = value.Name
	default:
		name = &ast.Ident{token.NoPos, "NoName", nil}
		p.addError("expected identifier or assignment")
	}
	if value == nil || p.tok == token.IDENT {
		vtype = p.parseIdent()
	}
	end := p.expect(token.RPAREN)

	ob := &ast.Object{
		NamePos: name.NamePos,
		Name:    name.Name,
		Kind:    ast.Var,
		Type:    vtype,
		Value:   value,
	}

	if old := p.curScope.Insert(ob); old != nil {
		p.addError("redeclaration of variable not allowed; original "+
			"declaration at: ", p.file.Position(old.NamePos))
	}

	return &ast.VarExpr{
		Expression: ast.Expression{Opening: open, Closing: end},
		Var:        varpos,
		Name:       name,
		Object:     ob,
	}
}

func (p *parser) tryExprOrList() ast.Expr {
	p.listok = true
	return p.parseGenExpr()
}
