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
		stmt, err := self.ParseStatement()
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		if stmt == nil {
			return statements, nil
		}
		statements = append(statements, stmt)
	}
}

func (self *Parser) ParseStatement() (*ast.Statement, error) {
	stmt := &ast.Statement{}
	peek, err := self.lexer.Peek()
	if err != nil {
		return nil, err
	}
	switch peek.Type {
	case token.Word:
		switch string(peek.Content) {
		case FunctionKeyWord:
			fun, err := self.ParseFunction()
			if err != nil {
				return nil, err
			}
			stmt.Node = fun
			return stmt, nil
		}
	case token.EOF:
		return nil, nil
	}
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
	case token.Integer:
		next, err := self.lexer.Peek()
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		switch next.Type {
		case token.Plus:
		case token.Minus:
		case token.Mult:
		case token.Div:
		default:
			return &ast.Expression{
				Node: &ast.IntegerLiteral{
					Value: string(tok.Content),
				},
			}, nil
		}
		// consume sign
		_, err = self.lexer.GetToken()
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		right, err := self.ParseExpr()
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		switch next.Type {
		case token.Plus:
			return &ast.Expression{
				Node: &ast.BinaryAdd{
					Left: &ast.Expression{
						Node: &ast.IntegerLiteral{
							Value: string(tok.Content),
						},
					},
					Right: right,
				},
			}, nil
		case token.Minus:
			return &ast.Expression{
				Node: &ast.BinaryMinus{
					Left: &ast.Expression{
						Node: &ast.IntegerLiteral{
							Value: string(tok.Content),
						},
					},
					Right: right,
				},
			}, nil
		case token.Mult:
			return &ast.Expression{
				Node: &ast.BinaryMult{
					Left: &ast.Expression{
						Node: &ast.IntegerLiteral{
							Value: string(tok.Content),
						},
					},
					Right: right,
				},
			}, nil
		case token.Div:
			return &ast.Expression{
				Node: &ast.BinaryDiv{
					Left: &ast.Expression{
						Node: &ast.IntegerLiteral{
							Value: string(tok.Content),
						},
					},
					Right: right,
				},
			}, nil
		}
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

func (self *Parser) ParseFunction() (*ast.FunctionDefinition, error) {
	// fun
	err := self.lexer.ExpectContent(
		token.Token{
			Type:    token.Word,
			Content: []rune("fun"),
		},
	)
	if err != nil {
		return nil, err
	}
	// name
	name, err := self.lexer.GetToken()
	if err != nil {
		return nil, err
	}
	if name.Type != token.Word {
		return nil, fmt.Errorf("unexpected token: %s", name)
	}
	// (
	err = self.lexer.ExpectType(
		token.Token{
			Type: token.OpenBracket,
		},
	)
	if err != nil {
		return nil, err
	}
	// TODO: no args for now
	// )
	err = self.lexer.ExpectType(
		token.Token{
			Type: token.CloseBracket,
		},
	)
	if err != nil {
		return nil, err
	}
	// {
	err = self.lexer.ExpectType(
		token.Token{
			Type: token.FigureOpenBracket,
		},
	)
	if err != nil {
		return nil, err
	}

	fun := &ast.FunctionDefinition{
		Name: ast.StringLiteral{
			Value: string(name.Content),
		},
	}
	for {
		peek, err := self.lexer.Peek()
		if err != nil {
			return nil, err
		}
		if peek.Type == token.FigureCloseBracket {
			_, err := self.lexer.GetToken()
			if err != nil {
				return nil, err
			}
			return fun, nil
		}
		stmt, err := self.ParseStatement()
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			return fun, nil
		}
		fun.Body = append(fun.Body, stmt)
	}
}
