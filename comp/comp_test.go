package comp_test

import (
  "testing"

  "github.com/rthornton128/calc1/comp"
)

func TestCompileFile(t *testing.T) {
  comp.CompileFile("c1", "(+ 5 3)")
  comp.CompileFile("c2", "(- 5 3)")
  comp.CompileFile("c3", "(- 10 2 4)")
  comp.CompileFile("c4", "(+ (+ 2 3) 3)")
  comp.CompileFile("c5", "(% (+ 2 8) 2)")
  //comp.CompileFile("test3", "(/ (* 4 3) (+ 3 3)")
}
