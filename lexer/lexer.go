//go:generate ragel -Z -G2 next.rl

package lexer

import (
	"errors"
	"fmt"
	"io"

	"github.com/db47h/monkey/token"
)

const bufferSize = 10

// queue is a FIFO queue.
//
type queue struct {
	items []token.Token
	head  int
	tail  int
	count int
}

func (q *queue) push(i token.Token) {
	if q.items == nil {
		q.items = make([]token.Token, 2)
	}
	if q.head == q.tail && q.count > 0 {
		items := make([]token.Token, len(q.items)*2)
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
func (q *queue) pop() token.Token {
	i := q.head
	q.head = (q.head + 1) % len(q.items)
	q.count--
	return q.items[i]
}

// Lexer wraps a lexer's state
//
type Lexer struct {
	queue
	r    io.Reader
	line int
	pos  int

	Copies int
	Reads  int

	// FSM state
	cs, ts, te, act int
	data            []byte
	bufSz           int
}

// New returns a new lexer for the given io.Reader.
//
func New(r io.Reader) *Lexer {
	l := &Lexer{
		r: r,
	}
	l.init()
	return l
}

// Reset resets the lexer to the start of the input stream. It will panic if
// the source reader is not an io.Seeker.
//
func (l *Lexer) Reset() {
	l.queue.count = 0
	l.queue.head = 0
	l.queue.tail = 0
	l.line = 0
	l.pos = 0
	l.Copies = 0
	l.Reads = 0
	l.bufSz = 0

	r := l.r.(io.Seeker)
	r.Seek(0, io.SeekStart)
	l.init()
}

func (l *Lexer) emit(pos int, typ token.Type, value interface{}) {
	l.push(token.Token{Pos: token.Pos(l.pos + pos), Type: typ, Literal: value})
}

func (l *Lexer) tokenString() string {
	return string(l.data[l.ts:l.te])
}

func (l *Lexer) newline(pos int) {
	l.line++
}

func (l *Lexer) updateBuffer() (p, pe, eof int, err error) {
	var n int
	if l.data == nil {
		l.data = make([]byte, bufferSize)
	}
	if l.bufSz >= len(l.data) {
		return l.bufSz, l.bufSz, l.bufSz, errors.New("buffer overrun")
	}
	for {
		n, err = l.r.Read(l.data[l.bufSz:])
		l.Reads++
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
		l.Copies++
		l.pos += l.ts
		l.bufSz = pe - l.ts
		copy(l.data[:], l.data[l.ts:pe])
		l.te -= l.ts
		l.ts = 0
	}
}

func (l *Lexer) QSize() int {
	return len(l.items)
}
