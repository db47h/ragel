// Copyright 2018 Denis Bernard <db047h@gmail.com>
// Licensed under the MIT license. See license text in the LICENSE file.

// Package ragel provides a ragel based scanner for Go.
//
package ragel

import (
	"errors"
	"fmt"
	"io"
)

const bufferSize = 32768

// A FSM provides an interface to ragel generated code.
//
type FSM interface {
	Init(l *Lexer)
	Run(l *Lexer, p, pe, eof int) (np, npe int)
	ErrState() int
}

// queue is a FIFO queue.
//
type queue struct {
	items []Token
	head  int
	tail  int
	count int
}

func (q *queue) push(i Token) {
	if q.items == nil {
		q.items = make([]Token, 2)
	}
	if q.head == q.tail && q.count > 0 {
		items := make([]Token, len(q.items)*2)
		copy(items, q.items[q.head:])
		copy(items[len(q.items)-q.head:], q.items[:q.head])
		q.head = 0
		q.tail = len(q.items)
		q.items = items
	}
	q.items[q.tail] = i
	q.tail = (q.tail + 1) % len(q.items)
	q.count++
}

// pop pops the first item from the queue. Callers must check that q.count > 0 beforehand.
//
func (q *queue) pop() Token {
	i := q.head
	q.head = (q.head + 1) % len(q.items)
	q.count--
	return q.items[i]
}

// Type represents the type of a token.
//
type Type int

// Token types.
//
const (
	Error Type = -1 + iota
	EOF
)

// Token represents a token.
//
type Token struct {
	Offset  int
	Type    Type
	Literal interface{}
}

// Lexer wraps a lexer's state
//
type Lexer struct {
	queue
	fileName string
	r        io.Reader
	offset   int
	line     int
	lines    []int

	// FSM state
	f               FSM
	cs, ts, te, act int
	data            []byte
	sz              int // size of unprocessed data
}

// New returns a new lexer for the given io.Reader.
//
func New(fileName string, r io.Reader, f FSM) *Lexer {
	l := &Lexer{
		r:        r,
		fileName: fileName,
		data:     make([]byte, bufferSize),
		lines:    []int{0},
		f:        f,
	}
	f.Init(l)
	return l
}

// Reset resets the lexer to the start of the input stream. It will panic if
// the source reader is not an io.Seeker.
//
func (l *Lexer) Reset() {
	l.queue.count = 0
	l.queue.head = 0
	l.queue.tail = 0
	l.offset = 0
	l.line = 0
	l.lines = l.lines[:1]
	l.sz = 0

	r := l.r.(io.Seeker)
	r.Seek(0, io.SeekStart)
	l.f.Init(l)
}

// Next returns the next token in the input stream.
//
func (l *Lexer) Next() Token {
	for l.count == 0 {
		var (
			n, p, pe, eof int
			err           error
		)

		if l.sz < len(l.data) {
			p = l.sz
			for {
				n, err = l.r.Read(l.data[p:])
				if n != 0 || err != nil {
					break
				}
			}
			pe = p + n
			eof = -1
			if n == 0 {
				eof = pe
			}
		} else {
			p, pe, eof = l.sz, l.sz, l.sz
			err = errors.New("buffer overrun")
		}

		p, pe = l.f.Run(l, p, pe, eof)

		if l.cs == l.f.ErrState() {
			// TODO: this needs to be a rune, and error acted upon.
			l.Emit(p, Error, fmt.Errorf("unexpected character %#U", l.data[p]))
		}

		if l.ts == 0 {
			l.sz = 0
			l.offset += p
		} else {
			l.offset += l.ts
			l.sz = pe - l.ts
			copy(l.data[:], l.data[l.ts:pe])
			l.te -= l.ts
			l.ts = 0
		}

		switch err {
		case nil:
		case io.EOF:
			l.Emit(0, EOF, nil)
		default:
			l.Emit(0, Error, err)
		}
	}

	return l.pop()
}

// SetState saves the fsm state. This function must be called before returning
// from FSM.Run and FSM.Init.
//
func (l *Lexer) SetState(cs, ts, te, act int) {
	l.cs, l.ts, l.te, l.act = cs, ts, te, act
}

// GetState gets the fsm state and buffer. This function must be called at
// the beginning of FSM.Run.
//
func (l *Lexer) GetState() (cs, ts, te, act int, data []byte) {
	return l.cs, l.ts, l.te, l.act, l.data
}

// Emit adds the given token to the output queue. The ts argument is usually
// the ragel variable ts.
//
func (l *Lexer) Emit(ts int, typ Type, value interface{}) {
	l.push(Token{Offset: l.offset + ts, Type: typ, Literal: value})
}

// Newline increments the line count. The p argument should be the ragel
// varaiable p.
//
func (l *Lexer) Newline(p int) {
	l.lines = append(l.lines, l.offset+p+1)
	l.line++
}

// Pos returns the line/column Position for the given offset.
//
func (l *Lexer) Pos(offset int) Position {
	i, j := 0, len(l.lines)
	for i < j {
		h := int(uint(i+j) >> 1)
		if !(l.lines[h] > offset) {
			i = h + 1
		} else {
			j = h
		}
	}
	return Position{l.fileName, i, offset - l.lines[i-1] + 1}
}

// Position represents the position of a token as its 1-based line and column numbers.
//
type Position struct {
	FileName string
	Line     int // 1-based line number
	Column   int // 1-based column number in bytes
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.FileName, p.Line, p.Column)
}
