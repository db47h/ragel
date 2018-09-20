// Copyright 2018 Denis Bernard <db047h@gmail.com>
// Licensed under the MIT license. See license text in the LICENSE file.

// Package ragel provides a driver for ragel based scanners.
//
// It provides all the functionality needed to read data from an  io.Reader,
// buffering, token position tracking (line and column) and error handling.
//
// The scanner reads input data in chunks of 32KiB (the default setting), then
// tokenizes the whole buffer and pushes the tokens found in a FIFO queue.
//
// scanner.Next pops tokens from the queue and populates it as needed.
//
// Error handling is done by returning a ragel.Error token. This can happen in 4
// cases:
//
//	- an IO error while reading from the input reader
//	- the input buffer gets full
//	- an illegal symbol is encountered (i.e. unhandled by the state machine definition)
//	- a ragel action calls State.Errorf
//
// The first two cases are non-recoverable errors; any subsequent call to
// scanner.Next returns ragel.EOF. For illegal symbols, scanning resumes at the
// following byte. Actions calling State.Errorf must update ragel's p variable
// (or issue an appropriate fgoto) for non-fatal errors.
//
// UTF8 input is supported via ragel's contib/unicode2ragel.rb.
// See lang_test.rl for an example use.
//
package ragel

import (
	"fmt"
	"io"
	"unicode/utf8"
)

// DefaultBufferSize is the default buffer size for a new Scanner. This can be
// changed at creation time with the BufferSize option.
//
const DefaultBufferSize = 32768

// Interface provides an interface to generated code for a given ragel state
// machine.
//
// The following canonical implementation must be present it the ragel state
// machine definition file:
//
//	type iface struct {}
//
//	func (iface) Init(s *ragel.State) (int, int) {
//		var cs, ts, te, act int
//		%%write init;
//		s.SaveVars(cs, ts, te, act)
//		return %%{ write start; }%%, %%{ write error; }%%
//	}
//
//	func (iface) Run(s *ragel.State, p, pe, eof int) (int, int) {
//		cs, ts, te, act, data := s.GetVars()
//		%%write exec;
//		s.SaveVars(cs, ts, te, act)
//		return p, pe
//	}
//
// The actual type name and whether it is public or private do not matter since
// it is only used by client code in the call to ragel.New:
//
//	scanner := ragel.New("stdin", os.Stdin, iface{})
//
type Interface interface {
	// Init initializes the FSM. Returns the start and error state indices.
	Init(s *State) (start, err int)
	// Run runs the scanner on the current buffer contents.
	Run(s *State, p, pe, eof int) (np, npe int)
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

// Option is a scanner option for New.
//
type Option interface {
	set(*state)
}

type optionFn func(*state)

func (f optionFn) set(s *state) { f(s) }

// BufferSize returns an Option that sets the input buffer size.
// The buffer size conditions the size of the longest possible token.
//
func BufferSize(size int) Option {
	return optionFn(func(s *state) {
		s.bufferSize = size
	})
}

// Scanner provides the API for scanner clients.
//
type Scanner state

// State holds a scanner's internal state while processing its input.
//
type State state

// state holds a scanner's internal state while processing its input.
//
type state struct {
	queue
	fileName   string
	r          io.Reader
	offset     int
	lines      []int
	iface      Interface
	bufferSize int

	// ragel state
	cs, ts, te, act int    // standard ragel variables
	data            []byte // standard ragel data ptr
	ss, se          int    // start and error state indices
	sz              int    // size of unprocessed data

	abort bool
}

// New returns a new scanner for the given io.Reader and Interface.
//
// The Interface implementation must be provided by the ragel generated code.
//
func New(fileName string, r io.Reader, iface Interface, opts ...Option) *Scanner {
	s := &state{
		r:          r,
		fileName:   fileName,
		lines:      []int{0},
		iface:      iface,
		bufferSize: DefaultBufferSize,
	}

	for _, o := range opts {
		o.set(s)
	}
	s.data = make([]byte, s.bufferSize)

	s.ss, s.se = iface.Init((*State)(s))
	return (*Scanner)(s)
}

// Reset resets the scanner to the start of the input stream. It will panic if
// the source reader is not an io.Seeker.
//
func (s *Scanner) Reset(opts ...Option) {
	s.queue.count = 0
	s.queue.head = 0
	s.queue.tail = 0
	s.abort = false
	s.offset = 0
	s.lines = s.lines[:1]
	s.sz = 0

	r := s.r.(io.Seeker)
	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	s.ss, s.se = s.iface.Init((*State)(s))

	for _, o := range opts {
		o.set((*state)(s))
	}
	if len(s.data) != s.bufferSize {
		s.data = make([]byte, s.bufferSize)
	}
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
			(*State)(s).Errorf(0, true, "token too long")
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
		if err != nil {
			eof = pe
		}

	again:
		p, pe = s.iface.Run((*State)(s), p, pe, eof)

		if s.cs == s.se {
			// error state: s.data[p:] did not match any rule. Report s.data[p]
			// as invalid and resume scanning past it.
			r, sz := utf8.DecodeRune(s.data[p:pe])
			if r == utf8.RuneError {
				r = rune(s.data[p])
			}
			(*State)(s).Errorf(p, false, "invalid character %#U", r)
			p += sz
			s.ts, s.te = p, p
			s.cs = s.ss
			if p < pe {
				goto again
			}
		}

		// iface.Run() sets ts to 0 in the start state and before checking if p == pe:
		//
		//	StartState:
		//	ts = 0
		//	if p == pe {
		//		return
		//	}
		//
		// as a result it is only safe to assume that the whole buffer has been
		// parsed and s.sz can only be set to 0 if s.cs = s.ss (start state).
		// Otherwise we just go on with the else clause which is a no-op if
		// pe == len(data) and will trigger the buffer overflow check at the top
		// of the loop.
		//
		if s.ts == 0 && s.cs == s.ss {
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
				(*State)(s).Errorf(0, true, err.Error())
			}
		}
	}

	if s.count > 0 {
		return s.pop()
	}
	return s.offset, EOF, "EOF"
}

// SaveVars saves the ragel state variables. These variables must be saved
// between calls to Interface.Run, so SaveVars must be called before returning
// from Interface.Run and Interface.Init.
//
func (s *State) SaveVars(cs, ts, te, act int) {
	s.cs, s.ts, s.te, s.act = cs, ts, te, act
}

// GetVars gets the ragel state variables and input buffer. This function must
// be called at the beginning of Interface.Run in order to initialize these
// variables properly.
//
func (s *State) GetVars() (cs, ts, te, act int, data []byte) {
	return s.cs, s.ts, s.te, s.act, s.data
}

// Emit adds the given token to the output queue. The ts argument is usually
// the ragel variable ts.
//
func (s *State) Emit(ts int, typ Token, value string) {
	s.push(item{off: s.offset + ts, tok: typ, lit: value})
}

// Errorf emits an Error token at position p. If the fatal argument is true,
// any subsequent call to Scanner.Next will return an EOF token, otherwise
// scanning will resume from position p (which the caller must update
// accordingly). The format and args arguments are passed through fmt.Sprintf to
// form the tokens's literal value.
//
func (s *State) Errorf(p int, fatal bool, format string, args ...interface{}) {
	if fatal {
		s.abort = true
	}
	s.push(item{off: s.offset + p, tok: Error, lit: fmt.Sprintf(format, args...)})
}

// Newline increments the line count. The p argument should be the ragel
// variable p.
//
func (s *State) Newline(p int) {
	s.lines = append(s.lines, s.offset+p+1)
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
