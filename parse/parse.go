package parse

import (
	"github.com/rthornton128/calc1/ast"
	"github.com/rthornton128/calc1/scan"
	"github.com/rthornton128/calc1/token"
)

func ParseFile(filename, src string) *ast.File {
	var p parser
	p.init(filename, src)
	return p.parseFile()
}

type parser struct {
	file    *token.File
	errors  scan.ErrorList
	scanner scan.Scanner

	pos token.Pos
	tok token.Token
	lit string
}

func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errors.Add(p.file.Position(pos), "'"+tok.String()+"'")
	}
	p.next()
	return pos
}

func (p *parser) init(fname, src string) {
	p.file = token.NewFile(fname, src)
	p.scanner.Init(p.file, src)
	p.next()
}

func (p *parser) next() {
	p.lit, p.tok, p.pos = p.scanner.Scan()
	if p.tok == token.COMMENT {
		p.next()
	}
}

func (p *parser) parseBasicLit() *ast.BasicLit {
	return &ast.BasicLit{LitPos: p.pos, Kind: p.tok, Lit: p.lit}
}

func (p *parser) parseBinaryExpr(open token.Pos) *ast.BinaryExpr {
	pos := p.pos
	op := p.tok
	p.next()

	var list []ast.Expr
	for p.tok != token.RPAREN {
		list = append(list, p.parseExpr())
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

func (p *parser) parseExpr() ast.Expr {
	var expr ast.Expr

	switch p.tok {
	case token.LPAREN:
		expr = p.parseExpression()
	case token.INTEGER:
		expr = p.parseBasicLit()
		p.next()
	default:
		p.errors.Add(p.file.Position(p.pos), "Expected expression, got"+
			p.tok.String())
	}

	return expr
}

func (p *parser) parseExpression() ast.Expr {
	var expr ast.Expr

	pos := p.expect(token.LPAREN)
	switch p.tok {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		expr = p.parseBinaryExpr(pos)
	default:
	}

	return expr
}

func (p *parser) parseFile() *ast.File {
	var expr ast.Expr
	expr = p.parseExpr()
	return &ast.File{Root: expr}
}
