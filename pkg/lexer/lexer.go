package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"

	"github.com/AvoidMe/weasel/pkg/token"
)

type Lexer struct {
	peeked *token.Token
	Body   *bufio.Reader
}

func New(body io.Reader) *Lexer {
	return &Lexer{
		Body: bufio.NewReader(body),
	}
}

func (self *Lexer) Peek() (*token.Token, error) {
	if self.peeked != nil {
		return self.peeked, nil
	}
	peeked, err := self.GetToken()
	if err != nil {
		return nil, fmt.Errorf("peek error: %w", err)
	}
	self.peeked = peeked
	return self.peeked, nil
}

func (self *Lexer) GetToken() (*token.Token, error) {
	if self.peeked != nil {
		peeked := self.peeked
		self.peeked = nil
		return peeked, nil
	}
	for {
		symbol, _, err := self.Body.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return &token.Token{Type: token.EOF}, nil
			}
			return nil, fmt.Errorf("read error: %w", err)
		}

		switch {
		case symbol == '(':
			return &token.Token{
				Type: token.OpenBracket,
			}, nil
		case symbol == ')':
			return &token.Token{
				Type: token.CloseBracket,
			}, nil
		case symbol == ',':
			return &token.Token{
				Type: token.Comma,
			}, nil
		case symbol == '"':
			err := self.Body.UnreadRune()
			if err != nil {
				return nil, fmt.Errorf("internal error: %w", err)
			}
			return self.ReadString()
		case symbol == '+':
			return &token.Token{
				Type: token.Plus,
			}, nil
		case symbol == '-':
			return &token.Token{
				Type: token.Minus,
			}, nil
		case symbol == '*':
			return &token.Token{
				Type: token.Mult,
			}, nil
		case symbol == '/':
			peek, err := self.Body.Peek(1)
			if err != nil {
				return nil, fmt.Errorf("internal error: %w", err)
			}
			if peek[0] != '/' {
				return &token.Token{
					Type: token.Div,
				}, nil
			}
			err = self.DiscardLine()
			if err != nil {
				return nil, fmt.Errorf("internal error: %w", err)
			}
		case symbol == '{':
			return &token.Token{
				Type: token.FigureOpenBracket,
			}, nil
		case symbol == '}':
			return &token.Token{
				Type: token.FigureCloseBracket,
			}, nil
		case unicode.IsDigit(symbol):
			err := self.Body.UnreadRune()
			if err != nil {
				return nil, fmt.Errorf("internal error: %w", err)
			}
			return self.ReadNumber()
		case unicode.IsSpace(symbol):
			continue
		case unicode.IsLetter(symbol):
			err := self.Body.UnreadRune()
			if err != nil {
				return nil, fmt.Errorf("internal error: %w", err)
			}
			return self.ReadWord()
		default:
			return nil, fmt.Errorf("unknown token: %s", string(symbol))
		}
	}
}

func (self *Lexer) ReadWord() (*token.Token, error) {
	var word []rune
	for {
		symbol, _, err := self.Body.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return &token.Token{
					Type:    token.Word,
					Content: word,
				}, nil
			}
			return nil, fmt.Errorf("read error: %w", err)
		}
		switch {
		case unicode.IsLetter(symbol), symbol == '_':
			word = append(word, symbol)
		case unicode.IsDigit(symbol):
			if len(word) == 0 {
				return nil, fmt.Errorf("internal error, expecting letter, got digit: %s", string(symbol))
			}
			word = append(word, symbol)
		default:
			err := self.Body.UnreadRune()
			if err != nil {
				return nil, fmt.Errorf("internal error: %w", err)
			}
			return &token.Token{
				Type:    token.Word,
				Content: word,
			}, nil
		}
	}
}

func (self *Lexer) ReadString() (*token.Token, error) {
	var word []rune
	for {
		symbol, _, err := self.Body.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("read error: %w", err)
		}
		switch symbol {
		// case '\n': // TODO: return parse error
		case '"':
			word = append(word, symbol)
			if len(word) == 1 {
				continue
			}
			return &token.Token{
				Type:    token.String,
				Content: word,
			}, nil
		default:
			word = append(word, symbol)
		}
	}
}

func (self *Lexer) ReadNumber() (*token.Token, error) {
	var number []rune
	for {
		symbol, _, err := self.Body.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return &token.Token{
					Type:    token.Integer,
					Content: number,
				}, nil
			}
			return nil, fmt.Errorf("read error: %w", err)
		}
		switch {
		case unicode.IsNumber(symbol):
			number = append(number, symbol)
		default:
			err := self.Body.UnreadRune()
			if err != nil {
				return nil, fmt.Errorf("unread error: %w", err)
			}
			return &token.Token{
				Type:    token.Integer,
				Content: number,
			}, nil
		}
	}
}

func (self *Lexer) DiscardLine() error {
	for {
		symbol, _, err := self.Body.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("read error: %w", err)
		}
		switch symbol {
		case '\n':
			return nil
		}
	}
}

func (self *Lexer) ExpectType(expect token.Token) error {
	tok, err := self.GetToken()
	if err != nil {
		return err
	}
	if tok.Type != expect.Type {
		return fmt.Errorf("unexpected token: %s, want: %s", tok, expect)
	}
	return nil
}

func (self *Lexer) ExpectContent(expect token.Token) error {
	tok, err := self.GetToken()
	if err != nil {
		return err
	}
	if tok.Type != expect.Type || string(tok.Content) != string(expect.Content) {
		return fmt.Errorf("unexpected token: %s, want: %s", tok, expect)
	}
	return nil
}
