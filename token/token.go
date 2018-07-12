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

// Pos represents the token's position in the input stream.
//
type Pos int

// A Token represents a Monkey token.
//
type Token struct {
	Pos     Pos
	Type    Type
	Literal interface{}
}
