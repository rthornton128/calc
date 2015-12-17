package cgen

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ir"
)

type cCompiler struct{ compiler }

/* Utility */

func cType(t ir.Type) string {
	switch t {
	case ir.Int:
		return "int32_t"
	case ir.Bool:
		return "bool"
	default:
		return "int"
	}
}

func (c *cCompiler) emitHeaders() {
	c.emitln("#include <stdio.h>")
	c.emitln("#include <stdint.h>")
	c.emitln("#include <stdbool.h>")
}

func (c *cCompiler) emitMain() {
	c.emitln("int main(void) {")
	c.emitln("printf(\"%d\\n\", _main());")
	c.emitln("return 0;")
	c.emitln("}")
}

/* Main Compiler */

func (c *cCompiler) compObject(o ir.Object) string {
	switch t := o.(type) {
	case *ir.Assignment:
		return c.compAssignment(t)
	case *ir.Constant:
		return c.compConstant(t)
	case *ir.Binary:
		return c.compBinary(t)
	case *ir.Call:
		return c.compCall(t)
	case *ir.For:
		return c.compFor(t)
	//case *ir.Function:
	//return c.compFunction(t)
	case *ir.If:
		return c.compIf(t)
	case *ir.Unary:
		return c.compUnary(t)
	case *ir.Var:
		return c.compVar(t)
	case *ir.Variable:
		c.emit("%s %s%d = 0;\n", cType(t.Type()), t.Name(), t.ID())
		return c.compVariable(t)
	}
	return ""
}

func (c *cCompiler) compAssignment(a *ir.Assignment) string {
	c.emit("%s%d = %s;\n", a.Lhs, a.Scope().Lookup(a.Lhs).ID(),
		c.compObject(a.Rhs))
	return a.Lhs
}

func (c *cCompiler) compBinary(b *ir.Binary) string {
	return fmt.Sprintf("(%s %s %s)",
		c.compObject(b.Lhs), b.Op.String(), c.compObject(b.Rhs))
}

func (c *cCompiler) compCall(call *ir.Call) string {
	args := make([]string, len(call.Args))
	for i, a := range call.Args {
		args[i] = fmt.Sprintf("%s", c.compObject(a))
	}
	return fmt.Sprintf("_%s(%s)", call.Name(), strings.Join(args, ","))
}

func (c *cCompiler) compConstant(con *ir.Constant) string {
	return con.String()
}

func (c *cCompiler) compDefine(d *ir.Define) string {
	switch t := d.Body.(type) {
	case *ir.Function:
		c.emit("%s {\n", c.compSignature(t))
		c.compFunction(d.Body.(*ir.Function))
		return ""
	default:
		return c.compObject(t)
	}
}

func (c *cCompiler) compFor(f *ir.For) string {
	c.emit("%s %s%d = 0;\n", cType(f.Type()), f.Name(), f.ID())
	c.emit("while (%s) {\n", c.compObject(f.Cond))
	for _, e := range f.Body[:len(f.Body)-1] {
		c.compObject(e)
	}
	c.emit("}\n%s%d = %s;\n", f.Name(), f.ID(),
		c.compObject(f.Body[len(f.Body)-1]))
	return fmt.Sprintf("%s%d", f.Name(), f.ID())
}

func (c *cCompiler) compFunction(f *ir.Function) {
	for _, e := range f.Body[:len(f.Body)-1] {
		c.compObject(e)
	}
	c.emit("return %s;\n}\n", c.compObject(f.Body[len(f.Body)-1]))
}

func (c *cCompiler) compIdent(i *ir.Var) string {
	return fmt.Sprintf("%s%d", i.Name(), i.Scope().Lookup(i.Name()).ID())
}

func (c *cCompiler) compIf(i *ir.If) string {
	c.emit("%s if%d = 0; /* %s */\n", cType(i.Type()), i.ID(), i.Name())
	c.emit("if (%s) {\n", c.compObject(i.Cond))
	c.emit("if%d = %s;\n", i.ID(), c.compObject(i.Then))
	if i.Else != nil {
		c.emitln("} else {")
		c.emit("if%d = %s;\n", i.ID(), c.compObject(i.Else))
	}
	c.emitln("}")
	return fmt.Sprintf("if%d", i.ID())
}

func (c *cCompiler) compPackage(p *ir.Package) {
	c.emitHeaders()
	names := p.Scope().Names()
	for _, name := range names {
		// later, this may need to check for import clauses
		if d, ok := p.Scope().Lookup(name).(*ir.Define); ok {
			if f, ok := d.Body.(*ir.Function); ok {
				c.emit("%s;\n", c.compSignature(f))
				params := make([]string, len(f.Params))
				for i, p := range f.Params {
					params[i] = cType(f.Scope().Lookup(p.Name()).(*ir.Param).Type())
				}
				c.emit("%s (*_%s)(%s) = f%d;\n", cType(f.Type()), d.Name(),
					strings.Join(params, ","), f.ID())
				defer c.compDefine(d)
			}
		}
	}
	c.emitMain()
}

func (c *cCompiler) compSignature(f *ir.Function) string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		param := f.Scope().Lookup(p.Name()).(*ir.Param)
		params[i] = fmt.Sprintf("%s %s%d", cType(param.Type()), param.Name(),
			param.ID())
	}
	return fmt.Sprintf("%s f%d(%s)", cType(f.Type()), f.ID(),
		strings.Join(params, ","))
}

func (c *cCompiler) compUnary(u *ir.Unary) string {
	return fmt.Sprintf("%s%s", u.Op, c.compObject(u.Rhs))
}

func (c *cCompiler) compVar(v *ir.Var) string {
	switch t := v.Scope().Lookup(v.Name()).(type) {
	case *ir.Define:
		return c.compDefine(t)
	case *ir.Param:
		return fmt.Sprintf("%s%d", t.Name(), t.ID())
	}
	panic("unreachable")
}

func (c *cCompiler) compVariable(v *ir.Variable) string {
	for _, p := range v.Params {
		param := v.Scope().Lookup(p.Name()).(*ir.Param)
		c.emit("%s %s%d = 0;\n", cType(param.Type()), param.Name(), param.ID())
	}
	for _, e := range v.Body[:len(v.Body)-1] {
		c.compObject(e)
	}
	c.emit("var%d = %s;\n", v.ID(), c.compObject(v.Body[len(v.Body)-1]))
	return fmt.Sprintf("var%d", v.ID())
}
