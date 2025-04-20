package ast

import (
	"fmt"
	"strings"
)

type Node interface {
	String() string
}

type Statements []*Statement

func (self Statements) String() string {
	var sb strings.Builder
	sb.WriteString("<ast.Statements: [")
	for i, n := range self {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(n.String())
	}
	sb.WriteString("]>")
	return sb.String()
}

type Statement struct {
	Node Node
}

func (self Statement) String() string {
	return fmt.Sprintf(
		"<ast.Statement: %s>",
		self.Node.String(),
	)
}

type Expression struct {
	Node Node
}

func (self Expression) String() string {
	return fmt.Sprintf(
		"<ast.Expression: %s>",
		self.Node.String(),
	)
}

type FunctionCall struct {
	Name StringLiteral
	Args []Node
}

func (self FunctionCall) String() string {
	var sb strings.Builder
	sb.WriteString(
		fmt.Sprintf(
			"<ast.FunctionCall (%s): [",
			self.Name,
		),
	)
	for i, n := range self.Args {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(n.String())
	}
	sb.WriteString("]>")
	return sb.String()
}

type StringLiteral struct {
	Value string
}

func (self StringLiteral) String() string {
	return fmt.Sprintf("<ast.String: %s>", self.Value)
}
