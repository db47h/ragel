//go:generate ragel -Z -G2 next.rl

package lexer

import (
	"bufio"
	"io"

	"github.com/db47h/monkey/token"
)

// Lexer wraps a lexer's state
//
type Lexer struct {
	r   *bufio.Reader
	c   chan token.Token
	pos int // char position of data buffer
}

// New returns a new lexer for the given io.Reader.
//
func New(r io.Reader) *Lexer {
	l := &Lexer{
		r: bufio.NewReader(r),
		c: make(chan token.Token),
	}
	l.init()
	return l
}

func (l *Lexer) emit(pos int, typ token.Type, value interface{}) {
	l.c <- token.Token{Pos: token.Pos(l.pos + pos), Type: typ, Literal: value}
}

// Next returns the next token in the input stream.
//
func (l *Lexer) Next() token.Token {
	return <-l.c
}
