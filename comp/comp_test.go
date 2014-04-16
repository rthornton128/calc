// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package comp_test

import (
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/rthornton128/calc/comp"
)

var ext string

func init() {
	ext = ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
}

func TestInteger(t *testing.T) {
	test_handler(t, "42", "42")
}

func TestSimpleExpression(t *testing.T) {
	test_handler(t, "(+ 5 3)", "8")
}

func TestSimpleExpressionWithComments(t *testing.T) {
	test_handler(t, ";comment 1\n(* 5 3); comment 2", "15")
}

func TestComplexExpression(t *testing.T) {
	test_handler(t, "(- (* 9 (+ 2 3)) (+ (/ 20 (% 15 10)) 1))", "40")
}

func test_handler(t *testing.T, src, expected string) {
	defer tearDown()

	comp.CompileFile("test", src)

	out, err := exec.Command("gcc"+ext, "-Wall", "-Wextra", "-std=c99",
		"--output=test"+ext, "test.c").CombinedOutput()

	if err != nil {
		t.Log(string(out))
		t.Fatal(err)
	}
	var output []byte

	switch runtime.GOOS {
	case "windows":
		output, err = exec.Command("test" + ext).Output()
	default:
		output, err = exec.Command("./test").Output()
	}

	if err != nil {
		t.Fatal(err)
	}
	if string(output) != expected {
		t.Fatal("For " + src + " expected " + expected + " got " + string(output))
	}
}

func tearDown() {
	os.Remove("test.c")
	os.Remove("test" + ext)
}
