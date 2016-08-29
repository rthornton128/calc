// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import (
	"io"

	"github.com/rthornton128/calc/ir"
)

// Registers and instructions specific to AMD64
func (c *X86) CGen(w io.Writer, pkg *ir.Package) {
	c.Writer = w
	//c.emit(".file %s\n", "xxx.calc")
	c.emit(".global main")
	for _, name := range pkg.Scope().Names() {
		if d, ok := pkg.Scope().Lookup(name).(*ir.Define); ok {
			if f, ok := d.Body.(*ir.Function); ok {
				c.emitf(".global _%s\n", name)
				defer func(name string) {
					c.emitf("_%s:", name)
					c.genObject(f, "%eax")
				}(name)
			}
		}
	}
	c.emit(".data")
	c.emitf("fmt: .asciz \"%%d\\12\"")
	c.emit("")
	c.emit(".text")
	c.emitMain()
}
