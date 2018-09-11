// Example CSV scanner with comma separated fields.

package main

import "github.com/db47h/ragel/v2"

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

type csv struct{}

// here we use a private implemetation for ragel.FSM since the scanner and
// parser are in the same package.
//
func (csv) Init(s *ragel.State) (int, int) {
	var cs, ts, te, act int
	%%write init;
	s.SaveVars(cs, ts, te, act)
	return %%{ write start; }%%, %%{ write error; }%%
}

func (csv) Run(s *ragel.State, p, pe, eof int) (int, int) {
	cs, ts, te, act, data := s.GetVars()
	%%write exec;
	s.SaveVars(cs, ts, te, act)
	return p, pe
}
