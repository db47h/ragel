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
	r      *bufio.Reader
	tokens []token.Token
	line   int // newline count
	pos    int // char position of data buffer

	cs, act, ts, te int
	data            []byte
	bufSz           int
	done            bool
}

// New returns a new lexer for the given io.Reader.
//
func New(r io.Reader) *Lexer {
	l := &Lexer{
		r:    bufio.NewReader(r),
		data: make([]byte, 10),
	}
	l.init()
	return l
}

func (l *Lexer) emit(pos int, typ token.Type, value interface{}) {
	l.tokens = append(l.tokens, token.Token{Pos: token.Pos(l.pos + pos), Type: typ, Literal: value})
}

func (l *Lexer) newline(pos int) {
	l.line++
}
