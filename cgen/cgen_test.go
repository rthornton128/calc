// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package comp_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/rthornton128/calc/comp"
)

var ext string

func init() {
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
}

func TestSimpleExpression(t *testing.T) {
	test_handler(t, "(decl main int 42)", "42")
}

func TestBinary(t *testing.T) {
	test_handler(t, "(decl main int (+ 5 3))", "8")
	test_handler(t, ";comment 1\n(decl main int (* 5 3)); comment 2", "15")
	test_handler(t, "(decl main int (- (* 9 (+ 2 3)) (+ (/ 20 (% 15 10)) 1)))",
		"40")
}

func TestFunc(t *testing.T) {
	test_handler(t, "(decl fn(a b int) int (+ a b))\n"+
		"(decl main int (fn 1 2))", "3")
}

func TestIfThenElse(t *testing.T) {
	test_handler(t, "(decl main int (if true int 99))", "99")
	test_handler(t, "(decl main int (if false int 2 3))", "3")
	test_handler(t, "(decl main int (if (< 2 3) int 7 3))", "7")
	test_handler(t, "(decl main int ((var a int) (if (< a 3) int 1 3)))", "1")
}

func TestVarAndAssign(t *testing.T) {
	test_handler(t, "(decl main int ((var (= a 5)) a))", "5")
	test_handler(t, "(decl main int ((var a int) (= a 42) a))", "42")
}

func TestUnary(t *testing.T) {
	test_handler(t, "(decl main int -24)", "-24")
	test_handler(t, "(decl main int ((var (= z 12)) -z))", "-12")
	test_handler(t, "(decl fn(num int) int -num)\n"+
		"(decl main int (fn -42))", "42")
}

func test_handler(t *testing.T, src, expected string) {
	defer tearDown()

	err := ioutil.WriteFile("test.calc", []byte(src), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = comp.CompileFile("test.calc", false)
	if err != nil {
		t.Log(src)
		t.Fatal(err)
	}
	os.Remove("test.calc")

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
	output = []byte(strings.TrimSpace(string(output)))

	if string(output) != expected {
		//t.Log("len output:", len(output))
		//t.Log("len expected:", len(expected))
		t.Fatal("For " + src + " expected " + expected + " got " + string(output))
	}
}

func tearDown() {
	os.Remove("test.c")
	os.Remove("test" + ext)
}
