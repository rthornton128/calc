package comp_test

import (
  "testing"

  "github.com/rthornton128/calc1/comp"
)

func TestCompileFile(t *testing.T) {
  comp.CompileFile("test1", "(+ (+ 2 3) 3)")
  comp.CompileFile("test2", "(% (+ 2 8) 2)")
  comp.CompileFile("test3", "(/ (* 4 3) (+ 3 3)")
}
