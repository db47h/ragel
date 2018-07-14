//go:generate ragel -Z -G2 next.rl

package lexer

import (
	"errors"
	"fmt"
	"io"

	"github.com/db47h/monkey/token"
)

const bufferSize = 32768

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

	// FSM state
	cs, ts, te, act int
	data            []byte
	sz              int // size of unprocessed data
}

// New returns a new lexer for the given io.Reader.
//
func New(r io.Reader) *Lexer {
	l := &Lexer{
		r:    r,
		data: make([]byte, bufferSize),
	}
	l.ragelInit()
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
	l.sz = 0

	r := l.r.(io.Seeker)
	r.Seek(0, io.SeekStart)
	l.ragelInit()
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

// Next returns the next token in the input stream.
//
func (l *Lexer) Next() token.Token {
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

		l.ragelNext(p, pe, eof)

		if l.cs == monkey_error {
			// TODO: this needs to be a rune
			l.emit(p, token.Error, fmt.Errorf("unexpected character %#U", l.data[p]))
		}

		if l.ts == 0 {
			l.sz = 0
			l.pos += p
		} else {
			l.pos += l.ts
			l.sz = pe - l.ts
			copy(l.data[:], l.data[l.ts:pe])
			l.te -= l.ts
			l.ts = 0
		}

		switch err {
		case nil:
		case io.EOF:
			l.emit(0, token.EOF, "EOF")
		default:
			l.emit(0, token.Error, err)
		}
	}

	return l.pop()
}
