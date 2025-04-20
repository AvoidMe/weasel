package parser

import (
	"fmt"

	"github.com/AvoidMe/weasel/pkg/ast"
	"github.com/AvoidMe/weasel/pkg/lexer"
	"github.com/AvoidMe/weasel/pkg/token"
)

type Parser struct {
	lexer *lexer.Lexer
}

func New(lexer *lexer.Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func Parse(lexer *lexer.Lexer) (ast.Statements, error) {
	return New(lexer).Parse()
}

func (self *Parser) Parse() (ast.Statements, error) {
	statements := ast.Statements{}
	for {
		tok, err := self.lexer.Peek()
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		if tok == nil {
			return statements, nil
		}
		stmt, err := self.ParseStatement()
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		statements = append(statements, stmt)
	}
}

func (self *Parser) ParseStatement() (*ast.Statement, error) {
	stmt := &ast.Statement{}
	expr, err := self.ParseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Node = expr
	return stmt, nil
}

func (self *Parser) ParseExpr() (*ast.Expression, error) {
	tok, err := self.lexer.GetToken()
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}
	switch tok.Type {
	case token.String:
		return &ast.Expression{
			Node: &ast.StringLiteral{
				Value: string(tok.Content),
			},
		}, nil
	case token.Word:
		peek, err := self.lexer.Peek()
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		switch peek.Type {
		case token.OpenBracket:
			// Parse function args
			call := &ast.FunctionCall{
				Name: ast.StringLiteral{Value: string(tok.Content)},
			}
			// consume "("
			_, err := self.lexer.GetToken()
			if err != nil {
				return nil, fmt.Errorf("parse error: %w", err)
			}
			// parse args
			for {
				peek, err := self.lexer.Peek()
				if err != nil {
					return nil, fmt.Errorf("parse error: %w", err)
				}
				if peek.Type == token.CloseBracket {
					// consume ")"
					_, err := self.lexer.GetToken()
					if err != nil {
						return nil, fmt.Errorf("parse error: %w", err)
					}
					return &ast.Expression{
						Node: call,
					}, nil
				}
				if peek.Type == token.Comma {
					// consume ","
					_, err := self.lexer.GetToken()
					if err != nil {
						return nil, fmt.Errorf("parse error: %w", err)
					}
				}
				expr, err := self.ParseExpr()
				if err != nil {
					return nil, fmt.Errorf("parse error: %w", err)
				}
				call.Args = append(call.Args, expr)
			}
		}
	default:
		return nil, fmt.Errorf("unexpected token: %s", tok.String())
	}
	panic("unreachable")
}
