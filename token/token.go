package token

// Type represents the type of a Monkey token.
//
type Type int

// Tokens
//
const (
	EOF Type = iota
	Error
	Ident
	Int
	Float
	Assign
	Plus
	Comma
	SemiColon
	LParen
	RParen
	LBrace
	RBrace
	Function
	Symbol
	Char
	String
)

// A Token represents a Monkey token.
//
type Token struct {
	Offset  int
	Type    Type
	Literal interface{}
}
