// Example CSV scanner with comma separated fields.
//
package csv

import "github.com/db47h/ragel"

%%{
	machine csv;

    newline = '\n' @{ s.Newline(p) };

	main := |*

    # comma
	',' {
		s.Emit(ts, Comma, string(data[ts:te]));
	};

	# Double Quote
	dliteralChar = newline | [^"] | '""' ;
	'"' . dliteralChar* . '"' {
        	s.Emit(ts, QuotedString, string(data[ts:te]))
	};

    # new lines
    newline {
        s.Emit(ts, EOL, "\n")
     };

    # numbers
    sign = [+\-];
	sign? digit+ ('.' digit+)* ('e' sign? digit+)? {
        	s.Emit(ts, Number, string(data[ts:te]))
	};

	# ignore whitespace
	(space - '\n')+;

    # do not exclude double quotes. This way the parser can check for bad quoting.
    (any - space - [,\n])+ { s.Emit(ts, String, string(data[ts:te])) };

	*|;
}%%

%%# anything beyond this point should be left unchanged

%%write data nofinal;

// here we use a private implemetation for ragel.FSM since the scanner and
// parser are in the same package.
//
type fsm struct {}

func (fsm) Init(s *ragel.Scanner) {
	var cs, ts, te, act int
	%%write init;
	s.SetState(cs, ts, te, act)
}

func (fsm) States() (start, err int) {
	return %%{ write start; }%%, %%{ write error; }%%
}

func (fsm) Run(s *ragel.Scanner, p, pe, eof int) (int, int) {
	cs, ts, te, act, data := s.GetState()
	%%write exec;
	s.SetState(cs, ts, te, act)
	return p, pe
}
