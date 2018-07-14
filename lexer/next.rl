package lexer

import (
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
        l.emit(l.ts, token.Ident, l.tokenString())
	};

	# Single Quote.
	sliteralChar = [^'\\] | newline | ( '\\' . any_count_line );
	'\'' . sliteralChar* . '\'' {
        l.emit(l.ts, token.Char, l.tokenString())
	};

	# Double Quote.
	dliteralChar = [^"\\] | newline | ( '\\' any_count_line );
	'"' . dliteralChar* . '"' {
        l.emit(l.ts, token.String, l.tokenString())
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
        l.emit(l.ts, token.Int, l.tokenString())
	};

	# Match a float. Upon entering the machine clear the buf, buffer
	# characters on every trans and dump the float upon leaving.
	digit+ '.' digit+ {
        l.emit(l.ts, token.Float, l.tokenString())
	};

	# Match a hex. Upon entering the hex part, clear the buf, buffer characters
	# on every trans and dump the hex on leaving transitions.
	'0x' xdigit+ {
        l.emit(l.ts, token.Int, l.tokenString())
	};

	*|;
}%%

%%write data nofinal;

func (l *Lexer) ragelInit() {
    %%write init;
}

func (l *Lexer) ragelNext(p, pe, eof int) (int, int) {
	%%write exec;
	return p, pe
}