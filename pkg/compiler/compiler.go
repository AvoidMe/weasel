package compiler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/AvoidMe/weasel/pkg/ast"
	"github.com/AvoidMe/weasel/pkg/lexer"
	"github.com/AvoidMe/weasel/pkg/parser"
)

type GoCompiler struct {
	Builder *strings.Builder
	Imports map[string]struct{}
}

func Compile(code []byte) ([]byte, error) {
	lexer := lexer.New(bytes.NewBuffer(code))
	ast, err := parser.Parse(lexer)
	if err != nil {
		return nil, err
	}
	compiler := GoCompiler{
		Builder: &strings.Builder{},
		Imports: make(map[string]struct{}),
	}
	return compiler.Compile(ast)
}

func (self *GoCompiler) Compile(body ast.Statements) ([]byte, error) {
	err := self.CompileStatements(body)
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	sb.WriteString("package main\n\n")

	sb.WriteString("import (\n")
	for k := range self.Imports {
		sb.WriteByte('"')
		sb.WriteString(k)
		sb.WriteByte('"')
		sb.WriteString("\n")
	}
	sb.WriteString(")\n\n")

	sb.WriteString("func main() {\n")
	sb.WriteString(self.Builder.String())
	sb.WriteString("}\n")

	return []byte(sb.String()), nil
}

func (self *GoCompiler) CompileStatements(body ast.Statements) error {
	for _, stmt := range body {
		switch v := stmt.Node.(type) {
		case *ast.FunctionDefinition:
			err := self.CompileFunctionDefinition(v)
			if err != nil {
				return err
			}
		case *ast.Expression:
			err := self.CompileExpression(v)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected statemnt: %s", stmt.Node)
		}
	}
	return nil
}

func (self *GoCompiler) CompileExpression(expr *ast.Expression) error {
	switch v := expr.Node.(type) {
	case *ast.StringLiteral:
		self.Builder.WriteString(v.Value)
	case *ast.IntegerLiteral:
		self.Builder.WriteString(v.Value)
	case *ast.BinaryAdd:
		self.CompileExpression(v.Left)
		self.Builder.WriteString(" + ")
		self.CompileExpression(v.Right)
	case *ast.BinaryMinus:
		self.CompileExpression(v.Left)
		self.Builder.WriteString(" - ")
		self.CompileExpression(v.Right)
	case *ast.BinaryMult:
		self.CompileExpression(v.Left)
		self.Builder.WriteString(" * ")
		self.CompileExpression(v.Right)
	case *ast.BinaryDiv:
		self.CompileExpression(v.Left)
		self.Builder.WriteString(" / ")
		self.CompileExpression(v.Right)
	case *ast.FunctionCall:
		switch v.Name.Value {
		case "print":
			self.Imports["fmt"] = struct{}{}
			self.Builder.WriteString("fmt.Println(")
		default:
			self.Builder.WriteString(v.Name.Value)
			self.Builder.WriteString("(")
		}
		for i, arg := range v.Args {
			switch v := arg.(type) {
			case *ast.Expression:
				if i > 0 {
					self.Builder.WriteString(", ")
				}
				err := self.CompileExpression(v)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unexpected func call node: %s", arg)
			}
		}
		self.Builder.WriteString(")\n")
	default:
		return fmt.Errorf("unexpected expr node: %s", expr.Node)
	}
	return nil
}

func (self *GoCompiler) CompileFunctionDefinition(def *ast.FunctionDefinition) error {
	self.Builder.WriteString(def.Name.Value)
	self.Builder.WriteString(" := func() {\n")
	err := self.CompileStatements(def.Body)
	if err != nil {
		return err
	}
	self.Builder.WriteString("\n}\n")
	return nil
}
