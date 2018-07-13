package lexer

import (
    "errors"
    "io"

    "github.com/db47h/monkey/token"
)

%%{
    machine monkey;

    # alphtype rune;
    access l.;

	newline = '\n' @{ l.newline(p) };
	any_count_line = any | newline;

	# Consume a C comment.
	c_comment := any_count_line* :>> '*/' @{fgoto main;};

	main := |*

	# Alpha numberic characters or underscore.
	alnum_u = alnum | '_';

	# Alpha charactres or underscore.
	alpha_u = alpha | '_';

	# Symbols. Upon entering clear the buffer. On all transitions
	# buffer a character. Upon leaving dump the symbol.
	( punct - [_'"] ) {
		l.emit(l.ts, token.Symbol, l.data[l.ts]);
	};

	# Identifier. Upon entering clear the buffer. On all transitions
	# buffer a character. Upon leaving, dump the identifier.
	alpha_u alnum_u* {
        l.emit(l.ts, token.Ident, string(l.data[l.ts:l.te]))
	};

	# Single Quote.
	sliteralChar = [^'\\] | newline | ( '\\' . any_count_line );
	'\'' . sliteralChar* . '\'' {
        l.emit(l.ts, token.Char, string(l.data[l.ts:l.te]))
	};

	# Double Quote.
	dliteralChar = [^"\\] | newline | ( '\\' any_count_line );
	'"' . dliteralChar* . '"' {
        l.emit(l.ts, token.String, string(l.data[l.ts:l.te]))
	};

	# Whitespace is standard ws, newlines and control codes.
	any_count_line - 0x21..0x7e;

	# Describe both c style comments and c++ style comments. The
	# priority bump on tne terminator of the comments brings us
	# out of the extend* which matches everything.
	'//' [^\n]* newline;

	'/*' { fgoto c_comment; };

	# Match an integer. We don't bother clearing the buf or filling it.
	# The float machine overlaps with int and it will do it.
	digit+ {
        l.emit(l.ts, token.Int, string(l.data[l.ts:l.te]))
	};

	# Match a float. Upon entering the machine clear the buf, buffer
	# characters on every trans and dump the float upon leaving.
	digit+ '.' digit+ {
        l.emit(l.ts, token.Float, string(l.data[l.ts:l.te]))
	};

	# Match a hex. Upon entering the hex part, clear the buf, buffer characters
	# on every trans and dump the hex on leaving transitions.
	'0x' xdigit+ {
        l.emit(l.ts, token.Int, string(l.data[l.ts:l.te]))
	};

	*|;
}%%

%%write data nofinal;

func (l *Lexer) init() {
    %%write init;
}


// Next returns the next token in the input stream.
//
func (l *Lexer) Next() token.Token {
    for !l.done && len(l.tokens) == 0 {
        var (
            p, pe, eof int
            n int
            err error
        )

        eof = -1
        p = l.bufSz

        if l.bufSz >= len(l.data) {
            l.done = true
            l.emit(p, token.Error, errors.New("buffer overrun"))
            break
        }
        for {
            n, err = l.r.Read(l.data[l.bufSz:])
            if n != 0 || err != nil {
                break
            }
        }
        pe = p + n

        if n == 0 {
            eof = pe
            l.done = true
            if err != nil && err != io.EOF {
                l.emit(p, token.Error, err)
                break
            }
        }

        %%write exec;

        if l.cs == monkey_error {
            l.done = true
            l.emit(p, token.Error, errors.New("parse error"))
            break
		}

        if l.ts == 0 {
            l.bufSz = 0
            l.pos += p
        } else {
            l.pos += l.ts
            l.bufSz = pe - l.ts
            copy(l.data, l.data[l.ts:pe])
            l.te -= l.ts
            l.ts = 0
        }
    }

    if len(l.tokens) == 0 {
        l.emit(0, token.EOF, "EOF")
    }

    t := l.tokens[0]
    copy(l.tokens, l.tokens[1:])
    l.tokens = l.tokens[:len(l.tokens)-1]
    return t
}