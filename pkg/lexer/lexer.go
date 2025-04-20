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
				return nil, nil
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
		case unicode.IsLetter(symbol):
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
