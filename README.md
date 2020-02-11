# ragel-go

[![godocb]][godoc]
[![goreportb]][goreport]

This package provides a driver for [ragel] based scanners in Go for streamed
input. It takes care of all the boilerplate code, letting the user focus on
the ragel state machine definition.

It is only intended for use with ragel machines defined as scanners. See
[section 6.3 of the ragel guide][ragel-guide].

## Rationale

Streaming input data to ragel -- i.e. reading from an io.Reader as opposed to
loading the whole input into memory, requires a *driver*, a piece of code that
buffers input and manipulates pointers for the ragel state machine.

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

In this example, the state machine actions just write tokens to stdout. This is
a lot of scaffolding and yet a few important things are still missing, like a
convenient way to retrieve the tokens and some mechanism to keep track of line
and column position for these tokens.

This package provides all the boilerplate code needed to read data from an
io.Reader, buffering, token position tracking (line and column) and error
handling and lets you focus on your ragel state machine definition.

## Quick start

In essence, creating a scanner is as simple as creating a ragel machine
definition for your language of choice, along with a tiny Go interface
implementation.

### Ragel state machine definition

Here is an abbreviated version for the C-like language mentioned above (with
support for comments, strings and char literals removed for brevity):

```go
// scanner/myc.rl

// Package scanner provides a scanner for my-C language.
//
package scanner // import "github.com/me/myc/scanner"

import "github.com/db47h/ragel/v2"

// Tokens for my-C
//
const (
    Ident ragel.Token = iota
    Float
    Symbol
)

// ragel state machine definition.
%%{
    machine lang;

    # utf-8 support
    include WChar "utf8.rl";

    newline = '\n' @{ s.Newline(p) };

    main := |*
        alpha_u = uletter | '_';
        alnum_u = alpha_u | digit;

        # punctuation
        punct {
            s.Emit(ts, Symbol, string(data[ts:te]));
        };

        # identifiers
        alpha_u alnum_u* {
            s.Emit(ts, Ident, string(data[ts:te]))
        };

        # ignore white space and control codes
        newline | 0x00..0x20 | 0x7f;

        # integers
        digit+ ('.' digit+)* {
            s.Emit(ts, Float, string(data[ts:te]))
        };
    *|;
}%%

%%write data nofinal;

// MyC implements ragel.Interface. This is the canonical implementation
// available at https://godoc.org/github.com/db47h/ragel#Interface
// with the type name changed to MyC.
//
type MyC struct {}

func (MyC) Init(s *ragel.State) (int, int) {
    var cs, ts, te, act int
    %%write init;
    s.SaveVars(cs, ts, te, act)
    return %%{ write start; }%%, %%{ write error; }%%
}

func (MyC) Run(s *ragel.State, p, pe, eof int) (int, int) {
    cs, ts, te, act, data := s.GetVars()
    %%write exec;
    s.SaveVars(cs, ts, te, act)
    return p, pe
}
```

Note that this example imports a file named "utf8.rl" for UTF-8 support. This
should be generated first and only once:

```bash
go run $GOPATH/github.com/db47h/ragel/cmd/unicode2ragel/main.go -o utf8.rl
```

Now supposing that the above example is saved in a file named `myc.rl`, run it
through ragel with the following command:

```bash
    ragel -Z -G2 myc.rl -o myc.go
```

This will generate a `myc.go` file. I tend to use the `-G2` because it yields
the fastest running code in most of my use cases. Your mileage may vary.

You might also want add `go:generate` directives in some other `.go` file in the
same package:

```go
//go:generate ragel -Z -G2 myc.rl -o myc.go
//go:generate sed -i "1s|.*|// Code generated by ragel. DO NOT EDIT.|" myc.go
```

The second directive changes the first line of the [generated code][codegen] so
that the generated file gets ignored by the Go tooling.

The `MyC` struct defined in the ragel state machine definition file is the
`ragel.Interface` implementation that is used in the call to `ragel.New`.
Besides changing the struct name to suit your needs (e.g. make it private), this
bit of code must be copy-pasted as-is from the [reference implementation][Interface].
If you feel that you need to change it for some reason, please file an issue.

### Scanning

Call `ragel.New` with an `io.Reader` and a `MyC` struct instance, then iterate
over tokens with `Scanner.Next`:

```go
// main.go

package main // import "github.com/me/myc"

import (
    "os"
    "fmt"

    "github.com/db47h/ragel/v2"
    "github.com/me/myc/scanner"
)

func main() {
    s := ragel.New("stdin", os.Stdin, scanner.MyC{})
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

## Supported Go and Ragel versions

### Go

This package uses the Go modules convention introduced in Go 1.11. The current
major version is v2, so it should be imported like this:

```go
    import "github.com/db47h/ragel/v2"
```

This is supported by Go 1.11, 1.10.4 and 1.9.7.

It might be possible to get it to work with earlier versions using third-party
package management tools, like dep, but this has not been tested.

### Ragel

Ragel provides good Go support from version 6.9 onwards.

The development releases (versions 7.x) should be avoided for the time being.
They have issues with that bit of code in the Go generator:

```go
    return %%{ write start; }%%, %%{ write error; }%%
```

If you need to use a 7.x version, just change the write statements like so:

```go
    return MachineName_start, MachineName_error
```

where `MachineName` is the FSM name as used in the ragel `machine` statement.

For all ragel versions, the alphabet data type needs to left to the default
setting of `byte`, even for UTF-8 input. i.e. if you have any `alphtype … ;`
statement, remove it.

## Handling UTF-8 input

UTF-8 input is supported by including UTF-8 state machines automatically
generated by the companion tool [unicode2ragel].

## Implementation details

The scanner reads input data in chunks of 32KiB (the default setting), then
tokenizes the whole buffer and pushes the tokens found in a FIFO queue.

scanner.Next pops tokens from the queue and populates it as needed.

An alternative to using a queue would be to run the scanning loop in a goroutine
and use a channel to emit tokens -- like text/template/parse in Go's standard
library. This approach is however 2 to 5 times slower than using a queue
(channels are great, just not here). See https://github.com/db47h/ragel/issues/4.

Error handling is done by returning a `ragel.Error` token. This can happen in
4 cases:

- an IO error while reading from the input reader
- the input buffer gets full
- an illegal symbol is encountered (i.e. unhandled by the state machine definition)
- a ragel action calls `State.Errorf`

The first two cases are non-recoverable errors; any subsequent call to
`scanner.Next` returns `ragel.EOF`. For illegal symbols, scanning resumes at the
following byte. Actions calling `State.Errorf` must update ragel's `p` variable
(or issue an appropriate `fgoto`) for non-fatal errors.

The standard ragel state variables `data`, `p`, `pe`, `eof`, `cs`, `ts`, `te`
and `act` are directly available in state machine definitions, i.e. `data[p]`
is the current character.

### The `ragel.Interface`

Ragel state machine definitions are templates where Go code is intermingled with
ragel directives (like the `%% write exec;` in the above example) and ragel
generates compilable Go code out of these templates.

The ragel directives `%% write exec;` and `%% write init;` generate code that
expects a few variables to be accessible in the program scope where these
directives appear. As a result, and although the variables names can be
customized (they could even be struct fields), the ragel state machine
definition is tightly coupled with the Go code that drives it. In order to have
a package that is go-gettable and updatable independently of users' own state
machine definitions, I needed to find a way to minimize the amount of Go code
present in a state machine definition file.

Since there is no way around some copy-paste, I came up with the [`ragel.Interface`][Interface]
interface and a [canonical implementation][Interface] that must be placed in
state machine definition files. This implementation uses as little code as
possible in order to limit potential bugs and avoid painful updates.

### API

The `ragel.Interface` implementation discussed above brings a `*State` instance
`s` into the scope of actions in ragel state machine definitions. Actions can
use its `Emit`, `Errorf` and `Newline` methods. The first two are mostly
self-explanatory. Users who need accurate line and column tracking must take
care of calling `Newline` at the right places (see how this is done in the above
example).

The `State.SaveVars` and `State.GetVars` methods manage the standard ragel state
variables and are only meant to be used by the `ragel.Interface` implementation.
As such, they must not be used by client code and should be considered private
for all intents and purposes.

The `Scanner` type wraps things together and provides the scanning client API.

Behind the scenes, `State` and `Scanner` are thin wrappers around `state`, their
sole purpose being a clean separation of concerns between the API used by state
machine definitions from the API used by parsers.

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
[Interface]: https://godoc.org/github.com/db47h/ragel#Interface
[godocb]: https://godoc.org/github.com/db47h/ragel?status.svg
[goreport]: https://goreportcard.com/report/github.com/db47h/ragel
[goreportb]: https://goreportcard.com/badge/github.com/db47h/ragel
[codegen]: https://golang.org/s/generatedcode
[ragel-guide]: http://www.colm.net/files/ragel/ragel-guide-6.10.pdf
[unicode2ragel]: https://godoc.org/github.com/db47h/ragel/cmd/unicode2ragel