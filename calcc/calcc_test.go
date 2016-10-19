package main_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

var tests = []struct {
	path string
	want int
}{
	{filepath.Join("tests", "min"), 0},
	{filepath.Join("tests", "binary"), 10},
	{filepath.Join("tests", "for"), 64},
	{filepath.Join("tests", "func"), 13},
	{filepath.Join("tests", "if"), 42},
	{filepath.Join("tests", "macros"), 7056},
	{filepath.Join("tests", "unary"), -24},
	{filepath.Join("tests", "var"), 8},
}

func TestASM(t *testing.T) {
	os := runtime.GOOS
	arch := runtime.GOARCH
	t.Logf("testing on %s with %s", os, arch)

	for _, test := range tests {
		t.Run(test.path,
			func(t *testing.T) {
				out, err := exec.Command("calcc", "-t", os, "-a", arch, "-opt=false",
					test.path+".calc").CombinedOutput()
				if err != nil {
					t.Errorf("test %s exited with error: %s", test.path, string(out))
				}
				out, err = exec.Command(test.path).Output()
				if err != nil {
					t.Errorf("test %s exited with error: %s", test.path, err)
				}
				re := regexp.MustCompile("\r?\n")
				out = re.ReplaceAll(out, []byte{})
				if string(out) != fmt.Sprint(test.want) {
					t.Errorf("want %d, got %s", test.want, out)
				}
			})
	}
}

func TestStandardC(t *testing.T) {
	for _, test := range tests {
		t.Run(test.path,
			func(t *testing.T) {
				out, err := exec.Command("calcc", "-opt=false",
					test.path+".calc").CombinedOutput()
				if err != nil {
					t.Errorf("test %s exited with error: %s", test.path, string(out))
				}
				out, err = exec.Command(test.path).Output()
				if err != nil {
					t.Errorf("test %s exited with error: %s", test.path, err)
				}
				re := regexp.MustCompile("\r?\n")
				out = re.ReplaceAll(out, []byte{})
				if string(out) != fmt.Sprint(test.want) {
					t.Errorf("want %d, got %s", test.want, out)
				}
			})
	}
}
