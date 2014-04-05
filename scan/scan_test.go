package scan_test

import (
  "testing"

  "github.com/rthornton128/calc1/scan"
  "github.com/rthornton128/calc1/token"
)

var src = "(+ 2 (- 4 1) (* 6 5) (% 10 2) (/ 9 3)); comment\n"

func TestScan(t *testing.T) {
  var expected = []token.Token{
    token.LPAREN,
    token.ADD,
    token.INTEGER,
    token.LPAREN,
    token.SUB,
    token.INTEGER,
    token.INTEGER,
    token.RPAREN,
    token.LPAREN,
    token.MUL,
    token.INTEGER,
    token.INTEGER,
    token.RPAREN,
    token.LPAREN,
    token.REM,
    token.INTEGER,
    token.INTEGER,
    token.RPAREN,
    token.LPAREN,
    token.QUO,
    token.INTEGER,
    token.INTEGER,
    token.RPAREN,
    token.RPAREN,
  }

  var s scan.Scanner
  s.Init(token.NewFile("", src), src)

  _, tok, pos := s.Scan()
  for i := 0; tok != token.EOF; i++ {
    if tok != expected[i] {
      t.Fatal(pos, "Expected:", expected[i], "Got:", tok)
    }
    _, tok, pos = s.Scan()
  }
}
