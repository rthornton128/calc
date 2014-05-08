// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package parse

import (
	"os"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/scan"
	"github.com/rthornton128/calc/token"
)

func ParseFile(filename, src string) *ast.File {
	var p parser
	p.init(filename, src)
	f := p.parseFile()
	if p.errors.Count() > 0 {
		p.errors.Print()
		return nil
	}
	return f
}

type parser struct {
	file    *token.File
	errors  scan.ErrorList
	scanner scan.Scanner

	curScope *ast.Scope
	topScope *ast.Scope

	pos token.Pos
	tok token.Token
	lit string
}

func (p *parser) addError(msg string) {
	p.errors.Add(p.file.Position(p.pos), msg)
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

func (p *parser) init(fname, src string) {
	p.file = token.NewFile(fname, 1, len(src))
	p.scanner.Init(p.file, src)
	p.curScope = ast.NewScope(nil)
	p.topScope = p.curScope
	p.next()
}

func (p *parser) next() {
	p.lit, p.tok, p.pos = p.scanner.Scan()
}

func (p *parser) openScope() {
	p.curScope = ast.NewScope(p.curScope)
}

func (p *parser) closeScope() {
	p.curScope = p.curScope.Parent()
}

func (p *parser) parseAssignExpr(open token.Pos) *ast.AssignExpr {
	pos := p.expect(token.ASSIGN)
	nam := p.parseIdent()
	val := p.parseGenExpr()

	ob := &ast.Object{
		NamePos: nam.NamePos,
		Name:    nam.Name,
		Type:    nil,
		Value:   val,
	}
	old := p.curScope.Lookup(nam.Name)
	if old == nil {
		p.addError("Cannot assign to undeclared identifier " + nam.Name)
	}
	if ob.Type.Name != old.Type.Name {
		p.addError("Cannot assign " + ob.Name + " of type (" + ob.Type.Name +
			") to " + old.Name + " of type (" + ob.Type.Name + ")")
	}
	end := p.expect(token.RPAREN)
	return &ast.AssignExpr{
		Expression: ast.Expression{Opening: open, Closing: end},
		Equal:      pos,
		Name:       nam,
		Object:     ob,
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
		p.addError("Expected '" + token.LPAREN.String() + "' or '" +
			token.INTEGER.String() + "' got '" + p.lit + "'")
		p.next()
	}

	return expr
}

func (p *parser) parseExpr() ast.Expr {
	var expr ast.Expr

	pos := p.expect(token.LPAREN)
	switch p.tok {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM,
		token.EQL, token.GTE, token.GTT, token.NEQ, token.LST, token.LTE:
		expr = p.parseBinaryExpr(pos)
	case token.VAR:
		expr = p.parseVarExpr(pos)
	default:
		p.addError("Expected binary operator but got '" + p.lit + "'")
	}

	return expr
}

func (p *parser) parseFile() *ast.File {
	var expr ast.Expr
	expr = p.parseGenExpr()
	if p.tok != token.EOF {
		p.addError("Expected EOF, got '" + p.lit + "'")
	}
	scope := ast.NewScope(nil)
	ob := &ast.Object{
		NamePos: token.NoPos,
		Name:    "main",
		Kind:    ast.Decl,
		Type:    &ast.Ident{NamePos: token.NoPos, Name: "int"},
		Value:   expr,
	}
	scope.Insert(ob)
	return &ast.File{Scope: scope}
}

func (p *parser) parseIdent() *ast.Ident {
	name := p.lit
	pos := p.expect(token.IDENT)
	return &ast.Ident{NamePos: pos, Name: name}
}

func (p *parser) parseVarExpr(lparen token.Pos) *ast.VarExpr {
	var (
		typ *ast.Ident
		val ast.Expr
	)
	varpos := p.expect(token.VAR)
	if p.tok != token.IDENT {
		p.expect(token.IDENT)
	}
	nam := p.parseIdent()

	// TODO: Needs improvement; maybe a p.tryTypeOrExpression?
	if p.tok == token.RPAREN {
		p.addError("Expected type, expression or literal, got: )")
	}

	if p.tok == token.IDENT {
		typ = p.parseIdent()
	}

	if p.tok != token.RPAREN {
		val = p.parseGenExpr()
	}

	if p.tok != token.RPAREN {
		typ = p.parseIdent()
	}
	// TODO: end
	rparen := p.expect(token.RPAREN)

	ob := &ast.Object{
		NamePos: nam.NamePos,
		Name:    nam.Name,
		Kind:    ast.Var,
		Type:    typ,
		Value:   val,
	}

	if old := p.curScope.Insert(ob); old != nil {
		p.addError("Identifier " + nam.Name + " redeclared; original " +
			"declaration at " + p.file.Position(old.NamePos).String())
	}

	return &ast.VarExpr{
		Expression: ast.Expression{Opening: lparen, Closing: rparen},
		Var:        varpos,
		Name:       nam,
		Object:     ob,
	}
}
