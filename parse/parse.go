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
	node := p.parseExpression()

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

func (p *parser) parseAssignExpr() *ast.AssignExpr {
	return &ast.AssignExpr{
		Equal: p.expect(token.ASSIGN),
		Name:  p.parseIdent(),
		Value: p.parseExpression(),
	}
}

func (p *parser) parseBasicLit() *ast.BasicLit {
	pos, tok, lit := p.pos, p.tok, p.lit
	p.next()
	return &ast.BasicLit{LitPos: pos, Kind: tok, Lit: lit}
}

func (p *parser) parseBinaryExpr() *ast.BinaryExpr {
	pos := p.pos
	op := p.tok
	p.next()

	return &ast.BinaryExpr{
		Op:    op,
		OpPos: pos,
		List:  p.parseExprList(),
	}
}

func (p *parser) parseCallExpr() *ast.CallExpr {
	return &ast.CallExpr{
		Name: p.parseIdent(),
		Args: p.parseExprList(),
	}
}

func (p *parser) parseDefineStmt() *ast.DefineStmt {
	p.expect(token.LPAREN)
	defer p.expect(token.RPAREN)

	return &ast.DefineStmt{
		Define: p.expect(token.DEFINE),
		Name:   p.parseIdent(),
		Type:   p.parseType(false),
		Body:   p.parseExpression(),
	}
}

func (p *parser) parseExpression() ast.Expr {
	var e ast.Expr
	switch p.tok {
	case token.LPAREN:
		p.expect(token.LPAREN)

		switch p.tok {
		case token.ADD, token.SUB, token.MUL, token.QUO, token.REM,
			token.EQL, token.GTE, token.GTT, token.NEQ, token.LST, token.LTE:
			e = p.parseBinaryExpr()
		case token.ASSIGN:
			e = p.parseAssignExpr()
		case token.FUNC:
			e = p.parseFuncExpr()
		case token.IDENT:
			e = p.parseCallExpr()
		case token.IF:
			e = p.parseIfExpr()
		case token.VAR:
			e = p.parseVarExpr()
		default:
			p.addError("Expected operator, keyword or identifier but got '" + p.lit +
				"'")
		}

		p.expect(token.RPAREN)
	case token.IDENT:
		e = p.parseIdent()
	case token.BOOL, token.INTEGER:
		e = p.parseBasicLit()
	case token.ADD, token.SUB:
		e = p.parseUnaryExpr()
	default:
		p.addError("Expected expression, got '" + p.lit + "'")
		p.next()
	}

	return e
}

func (p *parser) parseExprList() []ast.Expr {
	list := make([]ast.Expr, 0)
	for p.tok != token.RPAREN && p.tok != token.EOF {
		list = append(list, p.parseExpression())
	}
	return list
}

func (p *parser) parseFile() *ast.File {
	defs := make([]*ast.DefineStmt, 0)
	for p.tok != token.EOF {
		def := p.parseDefineStmt()

		switch def.Body.(type) {
		case *ast.FuncExpr:
			def.Kind = ast.FuncDecl
		default:
			def.Kind = ast.VarDecl
		}
		prev := p.curScope.Insert(&ast.Object{
			NamePos: def.Name.NamePos,
			Name:    def.Name.Name,
			Kind:    def.Kind,
		})
		if prev != nil {
			switch prev.Kind {
			case ast.FuncDecl:
				p.addError(prev.Name, " redeclared; declared as function at ",
					p.file.Position(prev.NamePos))
			case ast.VarDecl:
				p.addError(prev.Name, " redeclared; declared as variable at ",
					p.file.Position(prev.NamePos))
			}
			continue
		}

		defs = append(defs, def)
	}

	if len(defs) < 1 {
		p.addError("reached end of file without any declarations")
	}

	return &ast.File{Defs: defs}
}

func (p *parser) parseFuncExpr() *ast.FuncExpr {
	p.openScope()
	defer p.closeScope()

	return &ast.FuncExpr{
		Func:   p.expect(token.FUNC),
		Type:   p.parseType(true),
		Params: p.parseParamList(),
		Body:   p.parseExprList(),
	}
}

func (p *parser) parseIdent() *ast.Ident {
	name := p.lit
	return &ast.Ident{NamePos: p.expect(token.IDENT), Name: name}
}

func (p *parser) parseIfExpr() *ast.IfExpr {
	p.openScope()
	defer p.closeScope()

	ie := &ast.IfExpr{
		If:   p.expect(token.IF),
		Type: p.parseType(true),
		Cond: p.parseExpression(),
		Then: p.parseExpression(),
	}

	if p.tok != token.RPAREN {
		ie.Else = p.parseExpression()
	}
	return ie
}

func (p *parser) parseParamList() []*ast.Param {
	params := make([]*ast.Param, 0)
	p.expect(token.LPAREN)

	for p.tok != token.RPAREN {
		param := &ast.Param{Name: p.parseIdent(), Type: p.parseType(true)}
		o := &ast.Object{
			Kind:    ast.VarDecl,
			Name:    param.Name.Name,
			NamePos: param.Pos(),
		}
		if prev := p.curScope.Insert(o); prev != nil {
			p.addError("duplicate parameter ", prev.Name,
				"; previously declared at ", p.file.Position(prev.Pos()))
			continue
		}
		params = append(params, param)
	}
	p.expect(token.RPAREN)
	return params
}

func (p *parser) parseType(must bool) *ast.Ident {
	if !must && p.tok != token.COLON {
		return nil
	}
	p.expect(token.COLON)
	return p.parseIdent()
}

func (p *parser) parseUnaryExpr() *ast.UnaryExpr {
	pos, op := p.pos, p.lit
	p.next()
	return &ast.UnaryExpr{OpPos: pos, Op: op, Value: p.parseExpression()}
}

func (p *parser) parseVarExpr() *ast.VarExpr {
	return &ast.VarExpr{
		Var:    p.expect(token.VAR),
		Type:   p.parseType(true),
		Params: p.parseParamList(),
		Body:   p.parseExprList(),
	}
}
