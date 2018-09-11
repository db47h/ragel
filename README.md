# ragel-go

[![godocb]][godoc]
[![goreportb]][goreport]

This package provides a driver for [ragel] based scanners in Go for streamed
input. It takes care of all the boiler plate code, letting the user focus on
the ragel state machine definition.

## Intro

Streaming input data to ragel -- i.e. reading from an io.Reader as opposed to
loading the whole input into memory, requires a *driver*, a piece of code that
buffers input and manipulates pointers for the ragel engine.

Here is the Go version of the main loop from the C-like scanner in the ragel
examples:

```go
func scanner(r io.Reader) {
    var (
        buf = make([]byte, BUFSIZE)
        cs, act, ts, te int
        have, curLine = 0, 1
        done bool
    )

    // ragel init code
    %% write init;

    for !done {
        var (
            err error
            p = have
            ln, pe int
            eof = -1
        )

        if ( p == len(buf) ) {
            fmt.Fprintf(os.Stderr, "out of buffer space\n" )
            os.Exit(1)
        }

        for {
            ln, err = r.Read(buf[p:])
            if ln != 0 || err != nil {
                break
            }
        }
        pe = p + ln;

        // trigger end of file on any error
        if err != nil {
            eof = pe
            done = true
        }

        // Run the ragel state machine
        %% write exec;

        if cs == clang_error {
            fmt.Fprintf(os.Stderr, "parse error (invalid char)\n" )
            break
        }

        if ts == 0 {
            have = 0
        } else {
            // There is a prefix to preserve, shift it over.
            have = pe - ts;
            copy(buf, buf[ts:pe])
            te -= ts
            ts = 0
        }
    }
}
```

In this example, the state machine actions just write tokens to stdout. So we
have quite a lot of scaffolding and a few important things are missing, like a
convenient way to retrieve the tokens and some mechanism to keep track of line
and column position for these tokens.

## Implementation

The `Scanner` type provides all the functionality needed to read data from an
io.Reader, buffering, token position tracking (line and column) and error
handling.

The scanner reads input data in 32KiB chunks (the default buffer size), then
tokenizes the whole buffer and pushes the tokens found in a FIFO queue.

Note that this limits the size of the largest possible token to about 32767
bytes (and any token longer than this will not necessarily cause an error).

An alternative to using a queue would be to run the scanning loop in a goroutine
and use a channel to emit tokens (like text/template/parse in Go's standard
library). This approach is however slower than using a queue by an order of
magnitude.

Error handling is done by returning a `ragel.Error` token. This can happen in
3 cases: an IO error while reading from the input reader, the input buffer gets
full, or an illegal symbol is encountered. In the latter case, scanning
resumes at the following byte. The other two are non-recoverable errors; any
subsequent call to `scanner.Next` returns `ragel.EOF`.

Errors can also be generated from the state machine's actions, in which case the
caller decides if an error is fatal or not. For non-fatal errors, callers must
update ragel's `p` variable or issue an appropriate `fgoto`).

UTF8 input is supported via ragel's contib/unicode2ragel.rb.
See lang_test.rl for an example use.

## How to

In essence, creating a scanner is as simple as creating a ragel machine
definition for your language of choice, along with a tiny stub interface.

### Ragel state machine definition

Here is an abbreviated version for the C-like language mentioned above (with
support for comments, strings and char literals removed for brevity):

```go
// Package myc provides a scanner for my-C language.
package scan // import "github.com/me/myc/scan"

import "github.com/db47h/ragel"

// Tokens for my-C
//
const (
    Ident ragel.Token = iota
    Float
    Symbol
)

// state machine definition
%%{
    machine lang;
    include WChar "utf8.rl";

    newline = '\n' @{ s.Newline(p) };
    any_count_line = any | newline;

    main := |*

    alnum_u = ualnum | '_';

    alpha_u = ualpha | '_';

    punct {
        s.Emit(ts, Symbol, string(data[ts:te]));
    };

    alpha_u alnum_u* {
        s.Emit(ts, Ident, string(data[ts:te]))
    };

    newline | 0x00..0x20 | 0x7f;

    udigit+ ('.' udigit+)* {
        s.Emit(ts, Float, string(data[ts:te]))
    };
    *|;
}%%

%%write data nofinal;

// FSM implements ragel.FSM. This is the stub interface between the ragel
// generated code for a given machine and ragel.Scanner.
//
type FSM struct {}

func (FSM) Init(s *ragel.Scanner) (int, int) {
    var cs, ts, te, act int
    %%write init;
    s.SetState(cs, ts, te, act)
    return %%{ write start; }%%, %%{ write error; }%%
}

func (FSM) Run(s *ragel.Scanner, p, pe, eof int) (int, int) {
    cs, ts, te, act, data := s.GetState()
    %%write exec;
    s.SetState(cs, ts, te, act)
    return p, pe
}
```

Supposing that this is in a file named `myc.rl`, run it through ragel with the
following command:

```bash
    ragel -Z -G2 myc.rl
```

This will generate a `myc.go` file. I tend to use the `-G2` because it yields
the fastest running code in most of my use cases. Your mileage may vary.

You might also want add `go:generate` directives in some other `.go` file in the
same package:

```go
//go:generate ragel -Z -G2 myc.rl
//go:generate sed -i "1s|.*|// Code generated by ragel. DO NOT EDIT.|" myc.go
```

The second directive changes the first line of the [generated code][codegen] so that it
gets ignored by the Go tooling.

The `FSM` struct defined in the ragel file is the stub interface between the
generated code and `ragel.Scanner`. Besides changing the struct name to suit
your needs (e.g. make it private), this code should be used as-is. If you feel
that you need to change it for some reason, please file an issue.

### Scanning

Create a new `ragel.Scanner` and pass it the `FSM` struct defined above:

```go
package main

import (
    "os"
    "fmt"

    "github.com/db47h/ragel"
    "github.com/me/myc/scan"
)

func main() {
    s := ragel.New("stdin", os.Stdin, scan.FSM{})
    for {
        pos, tok, literal := s.Next()
        switch tok {
        case ragel.EOF:
            return
        case ragel.Error:
            fmt.Fprintf(os.Stderr, "scan error: %v: %v\n", s.Pos(pos), literal)
        default:
            fmt.Printf("%v: %d %s\n", s.Pos(), tok, literal)
        }
    }
}

```

That's it.

## API

The `Scanner` API has three users, each using a specific set of methods:

- the ragel state machine (actions depending on scanned input)
  - `Emit`
  - `Errorf`
  - `Newline`
- the stub interface
  - `SetState`
  - `GetState`
- the scanner client (code that wants to scan some input)
  - `New`
  - `Next`
  - `Pos`

This could be abstracted away using interfaces but it would affect performance
for no real benefit. Performance is the reason to use ragel in the first place.

## LICENSE

Copyright 2018 Denis Bernard <db047h@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[ragel]: http://www.colm.net/open-source/ragel/
[godoc]: https://godoc.org/github.com/db47h/ragel
[godocb]: https://godoc.org/github.com/db47h/ragel?status.svg
[goreport]: https://goreportcard.com/report/github.com/db47h/ragel
[goreportb]: https://goreportcard.com/badge/github.com/db47h/ragel
[codegen]: https://golang.org/s/generatedcode