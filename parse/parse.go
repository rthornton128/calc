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
	p.next()
}

func (p *parser) next() {
	p.lit, p.tok, p.pos = p.scanner.Scan()
}

func (p *parser) parseBasicLit() *ast.BasicLit {
	return &ast.BasicLit{LitPos: p.pos, Kind: p.tok, Lit: p.lit}
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
	case token.INTEGER:
		expr = p.parseBasicLit()
		p.next()
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
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		expr = p.parseBinaryExpr(pos)
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
	scope.Insert(&ast.Ident{NamePos: token.NoPos, Name: "main", Value: nil},
		expr)
	return &ast.File{Scope: scope}
}
