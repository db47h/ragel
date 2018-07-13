//go:generate ragel -Z -G2 next.rl

package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/db47h/monkey/token"
)

// Lexer wraps a lexer's state
//
type Lexer struct {
	r      *bufio.Reader
	tokens []token.Token
	line   int
	pos    int

	// FSM state
	cs, act, ts, te int
	data            []byte
	bufSz           int
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

func (l *Lexer) tokenString() string {
	return string(l.data[l.ts:l.te])
}

func (l *Lexer) newline(pos int) {
	l.line++
}

func (l *Lexer) updateBuffer() (p, pe, eof int, err error) {
	var n int
	if l.bufSz >= len(l.data) {
		return l.bufSz, l.bufSz, l.bufSz, errors.New("buffer overrun")
	}
	for {
		n, err = l.r.Read(l.data[l.bufSz:])
		if n != 0 || err != nil {
			break
		}
	}

	eof = -1
	if n == 0 {
		eof = pe
	}

	return l.bufSz, l.bufSz + n, eof, err
}

func (l *Lexer) shiftBuffer(p, pe int) {
	if l.cs == monkey_error {
		// TODO: this needs to be a rune
		l.emit(p, token.Error, fmt.Errorf("unexpected character %#U", l.data[p]))
	}

	if l.ts == 0 {
		l.bufSz = 0
		l.pos += p
	} else {
		l.pos += l.ts
		l.bufSz = pe - l.ts
		copy(l.data[:], l.data[l.ts:pe])
		l.te -= l.ts
		l.ts = 0
	}
}
