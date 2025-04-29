package token

import "fmt"

//go:generate stringer -type=TokenType
type TokenType int

const (
	Word TokenType = iota + 1
	String
	Integer
	Comma // ,

	Plus  // +
	Minus // -
	Mult  // *
	Div   // /

	OpenBracket        // (
	CloseBracket       // )
	FigureOpenBracket  // {
	FigureCloseBracket // }
	EOF
)

type Token struct {
	Type    TokenType
	Content []rune
}

func (self Token) String() string {
	if self.Content == nil {
		return fmt.Sprintf("<token=%s>", self.Type.String())
	}
	return fmt.Sprintf("<token=%s, value=%s>", self.Type.String(), string(self.Content))
}
