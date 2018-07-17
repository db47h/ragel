// Use this file as a template for your own scanner.
// Change the package name, custom token types and the ragel state machine definition.
// The FSM implementation at the bottom of the file should be left as-is.

package ragel_test

import "github.com/db47h/ragel"

// Custom token types.
//
const (
	Ident ragel.Token = iota
	Int
	Float
	Symbol
	Char
	String
)

%%{
	#state machine definition for your language

	machine lang;
	include WChar "utf8.rl";

	newline = '\n' @{ s.Newline(p) };
	any_count_line = any | newline;

	# Consume a C comment.
	c_comment := any_count_line* :>> '*/' @{fgoto main;};

	main := |*

	# Alpha numberic characters or underscore.
	alnum_u = ualnum | '_';

	# Alpha charactres or underscore.
	alpha_u = ualpha | '_';

	# Symbols. Upon entering clear the buffer. On all transitions
	# buffer a character. Upon leaving dump the symbol.
	( punct - [_'"] ) {
		s.Emit(ts, Symbol, string(data[ts:te]));
	};

	# Identifier. Upon entering clear the buffer. On all transitions
	# buffer a character. Upon leaving, dump the identifier.
	alpha_u alnum_u* {
        	s.Emit(ts, Ident, string(data[ts:te]))
	};

	# Single Quote.
	sliteralChar = [^'\\] | newline | ( '\\' . any_count_line );
	'\'' . sliteralChar* . '\'' {
        	s.Emit(ts, Char, string(data[ts:te]))
	};

	# Double Quote.
	dliteralChar = [^"\\] | newline | ( '\\' any_count_line );
	'"' . dliteralChar* . '"' {
        	s.Emit(ts, String, string(data[ts:te]))
	};

	# Whitespace is standard ws, newlines and control codes.
	newline | 0x00..0x20 | 0x7f;

	# Describe both c style comments and c++ style comments. The
	# priority bump on tne terminator of the comments brings us
	# out of the extend* which matches everything.
	'//' [^\n]* newline;

	'/*' { fgoto c_comment; };

	# Match an integer. We don't bother clearing the buf or filling it.
	# The float machine overlaps with int and it will do it.
	udigit+ {
        	s.Emit(ts, Int, string(data[ts:te]))
	};

	# Match a float. Upon entering the machine clear the buf, buffer
	# characters on every trans and dump the float upon leaving.
	udigit+ '.' udigit+ {
        	s.Emit(ts, Float, string(data[ts:te]))
	};

	# Match a hex. Upon entering the hex part, clear the buf, buffer characters
	# on every trans and dump the hex on leaving transitions.
	'0x' xdigit+ {
        	s.Emit(ts, Int, string(data[ts:te]))
	};

	*|;
}%%

%%# anything beyond this point should be left unchanged

%%write data nofinal;

// FSM implements ragel.FSM. This is the interface between the ragel generated
// code for a given machine and ragel.Scanner.
//
// Create a new scanner by calling scanner.New(..., fsm{}).
// 
type FSM struct {}

func (FSM) Init(s *ragel.Scanner) {
	var cs, ts, te, act int
	%%write init;
	s.SetState(cs, ts, te, act)
}

func (FSM) States() (start, err int) {
	return %%{ write start; }%%, %%{ write error; }%%
}

func (FSM) Run(s *ragel.Scanner, p, pe, eof int) (int, int) {
	cs, ts, te, act, data := s.GetState()
	%%write exec;
	s.SetState(cs, ts, te, act)
	return p, pe
}
