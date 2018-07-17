# ragel-go

This is a thin wrapper for [ragel] based scanners in Go. It takes care of all
the boiler plate code, letting the user focus on the ragel state machine
definition.

## Intro

This package is targetted at using the [ragel] state machine compiler for
lexical analysis of programming languages.

Besides writing the state machine definition, using ragel as a scanner requires
some scaffolding code to buffer input and track line and column position of
tokens read. Here is the Go version of the main loop from the C-like scanner in
ragel's examples:

```go
func scanner(r io.Reader) {
    var (
        buf = make([]byte, BUFSIZE)
        cs, act, ts, te int
        have, curline = 0, 1
        done bool
    )

    // ragel intit code
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

This is quite a lot of scaffolding and a few important things are missing, like
providing a `NextToken()` function: in this example, the state machine
actions just write tokens to stdout.

## Implementation

The `Scanner` type provides all the functionality needed to read data from an
io.Reader, buffering, token position tracking (line and column) and error
handling.

The scanner reads input data in 32KiB chunks (the default buffer size), then
tokenizes it and queues the tokens found. Scanner.Next takes care of this:
it populates the queue if its empty, otherwise it just pops tokens and
returns them.

Note that this limits the size of the largest possible token to about 32767
bytes.

An alternative to using a queue would be to run the scanning loop in a goroutine
and use a channel to emit tokens (like text/template/parse in Go's standard
library). This approach is however slower than using a queue by an order of
magnitude.

Error handling is done by returning a `ragel.Error` token. This can happen in
3 cases: an error while reading from the input reader, the input buffer gets
full, or an illegal symbol is encountered. In the latter case, scanning
resumes at the following byte. The other two are non-recoverable errors:
following these, Scanner.Next returns `ragel.EOF`.

Errors can also be generated from the state machine's actions, in which case
they are not considered fatal (but users must take care to update `p` or issue
an appropriate `fgoto`).

UTF8 input is supported via ragel's contib/unicode2ragel.rb.
See lang_test.rl for an example use.

## How to

In essence, the creating a scanner is as simple as creating a ragel machine
definition for your language of choice, along with a tiny stub interface.

### Ragel state machine definition

Here is an abbreviated version for the c-like language mentioned above (with
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

// CFSM implements ragel.FSM. This is the stub interface between the ragel
// generated code for a given machine and ragel.Scanner.
//
// Create a new scanner by calling scanner.New(..., fsm{}).
//
// Unless you have a very good reason, this code should be used as-is
//
type FSM struct {}

func (FSM) Init(s *ragel.Scanner) {
    var cs, ts, te, act int
    %%write init;
    s.SetState(cs, ts, te, act)
}

func (FSM) States() (start, err int) {
r   eturn %%{ write start; }%%, %%{ write error; }%%
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

This will generate a `myc.go` file.

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

[ragel]: http://www.colm.net/open-source/ragel/