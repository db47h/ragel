package lexer

import (
    "errors"
    "io"

    "github.com/db47h/monkey/token"
)

%%{
    machine monkey;

    # alphtype rune;
    # access l.;

	newline = '\n' @{curline += 1;};
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
		l.emit(ts, token.Symbol, data[ts]);
	};

	# Identifier. Upon entering clear the buffer. On all transitions
	# buffer a character. Upon leaving, dump the identifier.
	alpha_u alnum_u* {
        l.emit(ts, token.Ident, string(data[ts:te]))
	};

	# Single Quote.
	sliteralChar = [^'\\] | newline | ( '\\' . any_count_line );
	'\'' . sliteralChar* . '\'' {
        l.emit(ts, token.Char, string(data[ts:te]))
	};

	# Double Quote.
	dliteralChar = [^"\\] | newline | ( '\\' any_count_line );
	'"' . dliteralChar* . '"' {
        l.emit(ts, token.String, string(data[ts:te]))
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
        l.emit(ts, token.Int, string(data[ts:te]))
	};

	# Match a float. Upon entering the machine clear the buf, buffer
	# characters on every trans and dump the float upon leaving.
	digit+ '.' digit+ {
        l.emit(ts, token.Float, string(data[ts:te]))
	};

	# Match a hex. Upon entering the hex part, clear the buf, buffer characters
	# on every trans and dump the hex on leaving transitions.
	'0x' xdigit+ {
        l.emit(ts, token.Int, string(data[ts:te]))
	};

	*|;
}%%

%%write data nofinal;

func (l *Lexer) init() {
    go l.run()
}


func (l *Lexer) run()  {
    var cs, act, ts, te int

    %%write init;

    _ = act

    var have int
    var data [10]byte

    curline := 0
    done := false

    for !done {
        var p, pe, n, eof int
        var err error

        eof = -1
        p = have

        if have >= len(data) {
            l.emit(p, token.Error, errors.New("buffer overrun"))
            return
        }
        for {
            n, err = l.r.Read(data[have:])
            if n != 0 || err != nil {
                break
            }
        }
        pe = p + n

        if n == 0 {
            eof = pe
            done = true
            if err != nil && err != io.EOF {
                l.emit(p, token.Error, err)
            }
        }

        %%write exec;

        if cs == monkey_error {
            l.emit(p, token.Error, errors.New("parse error"))
            return
		}

        if ts == 0 {
            have = 0
            l.pos += p
        } else {
            l.pos += ts
            have = pe - ts
            copy(data[0:], data[ts:pe])
            te -= ts
            ts = 0
        }

    }
    l.emit(0, token.EOF, "EOF")
}