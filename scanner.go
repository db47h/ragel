// Copyright 2018 Denis Bernard <db047h@gmail.com>
// Licensed under the MIT license. See license text in the LICENSE file.

// Package ragel provides a thin wrapper for ragel based scanners.
//
// The Scanner type provides all the functionality needed to read data from an
// io.Reader, buffering, token position tracking (line and column) and error
// handling.
//
// The scanner reads input data in 32KiB chunks (the default buffer size), then
// tokenizes it and queues the tokens found. Scanner.Next takes care of this:
// it populates the queue if its empty, otherwise it just pops tokens and
// returns them.
//
// Note that this limits the size of the largest possible token to about 32767
// bytes.
//
// Error handling is done by returning a ragel.Error token. This can happen in
// 3 cases: an error while reading from the input reader, the input buffer gets
// full, or an illegal symbol is encountered. In the latter case, scanning
// resumes at the following byte. The other two are non-recoverable errors:
// following these, Scanner.Next returns ragel.EOF.
//
// UTF8 input is supported via ragel's contib/unicode2ragel.rb.
// See lang_test.rl for an example use.
//
package ragel // import "github.com/db47h/ragel"

import (
	"fmt"
	"io"
	"unicode/utf8"
)

const bufferSize = 32768

// A FSM provides an interface to ragel generated code for a given state
// machine. The implementation is part of the ragel machine definition file.
//
// See lang_test.rl for a typical implementation.
//
type FSM interface {
	Init(l *Scanner)
	States() (start, err int)
	Run(l *Scanner, p, pe, eof int) (np, npe int)
}

// queue is a FIFO queue.
//
type queue struct {
	items []item
	head  int
	tail  int
	count int
}

func (q *queue) push(i item) {
	if q.items == nil {
		q.items = make([]item, 2)
	}
	if q.head == q.tail && q.count > 0 {
		items := make([]item, len(q.items)*2)
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
func (q *queue) pop() (int, Token, string) {
	i := q.head
	q.head = (q.head + 1) % len(q.items)
	q.count--
	pi := &q.items[i]
	return pi.off, pi.tok, pi.lit
}

// Token represents a lexical token.
//
type Token int

// Token types.
//
const (
	Error Token = -2 + iota
	EOF
)

// An item wraps a token with its position and literal value.
//
type item struct {
	off int
	tok Token
	lit string
}

// Scanner holds the scanner's internal state while processing its input.
//
type Scanner struct {
	queue
	fileName string
	r        io.Reader
	offset   int
	line     int
	lines    []int
	abort    bool

	// FSM state
	f               FSM
	cs, ts, te, act int
	data            []byte
	sz              int // size of unprocessed data
}

// New returns a new scanner for the given io.Reader and FSM.
//
// The FSM implementation is provided by the ragel generated code.
//
func New(fileName string, r io.Reader, f FSM) *Scanner {
	s := &Scanner{
		r:        r,
		fileName: fileName,
		data:     make([]byte, bufferSize),
		lines:    []int{0},
		f:        f,
	}
	f.Init(s)
	return s
}

// Reset resets the scanner to the start of the input stream. It will panic if
// the source reader is not an io.Seeker.
//
func (s *Scanner) Reset() {
	s.queue.count = 0
	s.queue.head = 0
	s.queue.tail = 0
	s.abort = false
	s.offset = 0
	s.line = 0
	s.lines = s.lines[:1]
	s.sz = 0

	r := s.r.(io.Seeker)
	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	s.f.Init(s)
}

// Next returns the next token in the input stream.
//
func (s *Scanner) Next() (offset int, token Token, literal string) {
	for s.count == 0 && !s.abort {
		var (
			n, p, pe, eof int
			err           error
		)

		if s.sz >= len(s.data) {
			s.Errorf(0, "buffer overrun")
			s.abort = true
			break
		}

		p = s.sz
		for {
			n, err = s.r.Read(s.data[p:])
			if n != 0 || err != nil {
				break
			}
		}
		pe = p + n
		eof = -1
		if n == 0 {
			eof = pe
		}

	again:
		p, pe = s.f.Run(s, p, pe, eof)

		if ss, se := s.f.States(); s.cs == se {
			r, _ := utf8.DecodeRune(s.data[p:])
			if r == utf8.RuneError {
				r = rune(s.data[p])
			}
			s.Errorf(p, "invalid character %#U", r)
			p++
			s.cs = ss
			if p < pe {
				goto again
			}
		}

		if s.ts == 0 {
			s.sz = 0
			s.offset += p
		} else {
			s.offset += s.ts
			s.sz = pe - s.ts
			copy(s.data[:], s.data[s.ts:pe])
			s.te -= s.ts
			s.ts = 0
		}

		if err != nil {
			s.abort = true
			if err != io.EOF {
				s.Emit(0, Error, err.Error())
			}
		}
	}

	if s.count > 0 {
		return s.pop()
	}
	return s.offset, EOF, "EOF"
}

// SetState saves the fsm state. This function must be called before returning
// from FSM.Run and FSM.Init.
//
func (s *Scanner) SetState(cs, ts, te, act int) {
	s.cs, s.ts, s.te, s.act = cs, ts, te, act
}

// GetState gets the fsm state and buffer. This function must be called at
// the beginning of FSM.Run.
//
func (s *Scanner) GetState() (cs, ts, te, act int, data []byte) {
	return s.cs, s.ts, s.te, s.act, s.data
}

// Emit adds the given token to the output queue. The ts argument is usually
// the ragel variable ts.
//
func (s *Scanner) Emit(ts int, typ Token, value string) {
	s.push(item{off: s.offset + ts, tok: typ, lit: value})
}

// Errorf is a wrapper around Emit that emits error tokens. format and args
// are passed through fmt.Sprintf to forom the literal value.
//
func (s *Scanner) Errorf(ts int, format string, args ...interface{}) {
	s.push(item{off: s.offset + ts, tok: Error, lit: fmt.Sprintf(format, args...)})
}

// Newline increments the line count. The p argument should be the ragel
// varaiable p.
//
func (s *Scanner) Newline(p int) {
	s.lines = append(s.lines, s.offset+p+1)
	s.line++
}

// Pos returns the line/column Position for the given offset.
//
func (s *Scanner) Pos(offset int) Position {
	i, j := 0, len(s.lines)
	for i < j {
		h := int(uint(i+j) >> 1)
		if !(s.lines[h] > offset) {
			i = h + 1
		} else {
			j = h
		}
	}
	return Position{s.fileName, i, offset - s.lines[i-1] + 1}
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
