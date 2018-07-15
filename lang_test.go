
//line lang_test.rl:1
// Use this file as a template for your own scanner.
// Change the package name, custom token types and the ragel machine definition.
package ragel_test

import "github.com/db47h/ragel"

// Custom token types.
const (
	Ident ragel.Token = ragel.EOF + iota
	Int
	Float
	Symbol
	Char
	String
)

// The FSM definition of your language.

//line lang_test.rl:90


// anything beyond this point should be left unchanged.


//line lang_test.go:28
const lang_start int = 10
const lang_error int = 0

const lang_en_c_comment int = 8
const lang_en_main int = 10


//line lang_test.rl:95

type fsm struct {}

func (fsm) ErrState() int { return 0 }

func (fsm) Init(s *ragel.Scanner) {
	var cs, ts, te, act int
	
//line lang_test.go:45
	{
	cs = lang_start
	ts = 0
	te = 0
	act = 0
	}

//line lang_test.rl:103
	s.SetState(cs, ts, te, act)
}

func (fsm) Run(s *ragel.Scanner, p, pe, eof int) (int, int) {
	cs, ts, te, act, data := s.GetState()
	
//line lang_test.go:60
	{
	if p == pe {
		goto _test_eof
	}
	switch cs {
	case 10:
		goto st_case_10
	case 1:
		goto st_case_1
	case 2:
		goto st_case_2
	case 3:
		goto st_case_3
	case 4:
		goto st_case_4
	case 11:
		goto st_case_11
	case 5:
		goto st_case_5
	case 12:
		goto st_case_12
	case 6:
		goto st_case_6
	case 13:
		goto st_case_13
	case 14:
		goto st_case_14
	case 7:
		goto st_case_7
	case 15:
		goto st_case_15
	case 16:
		goto st_case_16
	case 8:
		goto st_case_8
	case 9:
		goto st_case_9
	case 17:
		goto st_case_17
	case 0:
		goto st_case_0
	}
	goto st_out
tr2:
//line lang_test.rl:57
te = p+1
{
        	s.Emit(ts, String, string(data[ts:te]))
	}
	goto st10
tr6:
//line lang_test.rl:51
te = p+1
{
        	s.Emit(ts, Char, string(data[ts:te]))
	}
	goto st10
tr8:
//line lang_test.rl:39
p = (te) - 1
{
		s.Emit(ts, Symbol, string(data[ts:te]));
	}
	goto st10
tr10:
//line lang_test.rl:23
 s.Newline(p) 
//line lang_test.rl:67
te = p+1

	goto st10
tr11:
//line lang_test.rl:73
p = (te) - 1
{
        	s.Emit(ts, Int, string(data[ts:te]))
	}
	goto st10
tr18:
//line lang_test.rl:62
te = p+1

	goto st10
tr19:
//line lang_test.rl:23
 s.Newline(p) 
//line lang_test.rl:62
te = p+1

	goto st10
tr20:
//line lang_test.rl:39
te = p+1
{
		s.Emit(ts, Symbol, string(data[ts:te]));
	}
	goto st10
tr25:
//line lang_test.rl:39
te = p
p--
{
		s.Emit(ts, Symbol, string(data[ts:te]));
	}
	goto st10
tr26:
//line lang_test.rl:69
te = p+1
{ {goto st8 } }
	goto st10
tr27:
//line lang_test.rl:73
te = p
p--
{
        	s.Emit(ts, Int, string(data[ts:te]))
	}
	goto st10
tr30:
//line lang_test.rl:79
te = p
p--
{
        	s.Emit(ts, Float, string(data[ts:te]))
	}
	goto st10
tr31:
//line lang_test.rl:85
te = p
p--
{
        	s.Emit(ts, Int, string(data[ts:te]))
	}
	goto st10
tr32:
//line lang_test.rl:45
te = p
p--
{
        	s.Emit(ts, Ident, string(data[ts:te]))
	}
	goto st10
	st10:
//line NONE:1
ts = 0

		if p++; p == pe {
			goto _test_eof10
		}
	st_case_10:
//line NONE:1
ts = p

//line lang_test.go:214
		switch data[p] {
		case 10:
			goto tr19
		case 34:
			goto st1
		case 39:
			goto st3
		case 47:
			goto tr21
		case 48:
			goto tr22
		case 95:
			goto st16
		}
		switch {
		case data[p] < 65:
			switch {
			case data[p] < 49:
				if 33 <= data[p] && data[p] <= 46 {
					goto tr20
				}
			case data[p] > 57:
				if 58 <= data[p] && data[p] <= 64 {
					goto tr20
				}
			default:
				goto tr23
			}
		case data[p] > 90:
			switch {
			case data[p] < 97:
				if 91 <= data[p] && data[p] <= 96 {
					goto tr20
				}
			case data[p] > 122:
				if 123 <= data[p] && data[p] <= 126 {
					goto tr20
				}
			default:
				goto st16
			}
		default:
			goto st16
		}
		goto tr18
tr1:
//line lang_test.rl:23
 s.Newline(p) 
	goto st1
	st1:
		if p++; p == pe {
			goto _test_eof1
		}
	st_case_1:
//line lang_test.go:269
		switch data[p] {
		case 10:
			goto tr1
		case 34:
			goto tr2
		case 92:
			goto st2
		}
		goto st1
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
		if data[p] == 10 {
			goto tr1
		}
		goto st1
tr5:
//line lang_test.rl:23
 s.Newline(p) 
	goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line lang_test.go:297
		switch data[p] {
		case 10:
			goto tr5
		case 39:
			goto tr6
		case 92:
			goto st4
		}
		goto st3
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
		if data[p] == 10 {
			goto tr5
		}
		goto st3
tr21:
//line NONE:1
te = p+1

	goto st11
	st11:
		if p++; p == pe {
			goto _test_eof11
		}
	st_case_11:
//line lang_test.go:326
		switch data[p] {
		case 42:
			goto tr26
		case 47:
			goto st5
		}
		goto tr25
	st5:
		if p++; p == pe {
			goto _test_eof5
		}
	st_case_5:
		if data[p] == 10 {
			goto tr10
		}
		goto st5
tr22:
//line NONE:1
te = p+1

	goto st12
	st12:
		if p++; p == pe {
			goto _test_eof12
		}
	st_case_12:
//line lang_test.go:353
		switch data[p] {
		case 46:
			goto st6
		case 120:
			goto st7
		}
		if 48 <= data[p] && data[p] <= 57 {
			goto tr23
		}
		goto tr27
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
		if 48 <= data[p] && data[p] <= 57 {
			goto st13
		}
		goto tr11
	st13:
		if p++; p == pe {
			goto _test_eof13
		}
	st_case_13:
		if 48 <= data[p] && data[p] <= 57 {
			goto st13
		}
		goto tr30
tr23:
//line NONE:1
te = p+1

	goto st14
	st14:
		if p++; p == pe {
			goto _test_eof14
		}
	st_case_14:
//line lang_test.go:392
		if data[p] == 46 {
			goto st6
		}
		if 48 <= data[p] && data[p] <= 57 {
			goto tr23
		}
		goto tr27
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st15
			}
		case data[p] > 70:
			if 97 <= data[p] && data[p] <= 102 {
				goto st15
			}
		default:
			goto st15
		}
		goto tr11
	st15:
		if p++; p == pe {
			goto _test_eof15
		}
	st_case_15:
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st15
			}
		case data[p] > 70:
			if 97 <= data[p] && data[p] <= 102 {
				goto st15
			}
		default:
			goto st15
		}
		goto tr31
	st16:
		if p++; p == pe {
			goto _test_eof16
		}
	st_case_16:
		if data[p] == 95 {
			goto st16
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st16
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st16
			}
		default:
			goto st16
		}
		goto tr32
tr15:
//line lang_test.rl:23
 s.Newline(p) 
	goto st8
	st8:
//line NONE:1
ts = 0

		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
//line lang_test.go:469
		switch data[p] {
		case 10:
			goto tr15
		case 42:
			goto st9
		}
		goto st8
	st9:
		if p++; p == pe {
			goto _test_eof9
		}
	st_case_9:
		switch data[p] {
		case 10:
			goto tr15
		case 42:
			goto st9
		case 47:
			goto tr17
		}
		goto st8
tr17:
//line lang_test.rl:27
{goto st10 }
	goto st17
	st17:
		if p++; p == pe {
			goto _test_eof17
		}
	st_case_17:
//line lang_test.go:500
		goto st0
st_case_0:
	st0:
		cs = 0
		goto _out
	st_out:
	_test_eof10: cs = 10; goto _test_eof
	_test_eof1: cs = 1; goto _test_eof
	_test_eof2: cs = 2; goto _test_eof
	_test_eof3: cs = 3; goto _test_eof
	_test_eof4: cs = 4; goto _test_eof
	_test_eof11: cs = 11; goto _test_eof
	_test_eof5: cs = 5; goto _test_eof
	_test_eof12: cs = 12; goto _test_eof
	_test_eof6: cs = 6; goto _test_eof
	_test_eof13: cs = 13; goto _test_eof
	_test_eof14: cs = 14; goto _test_eof
	_test_eof7: cs = 7; goto _test_eof
	_test_eof15: cs = 15; goto _test_eof
	_test_eof16: cs = 16; goto _test_eof
	_test_eof8: cs = 8; goto _test_eof
	_test_eof9: cs = 9; goto _test_eof
	_test_eof17: cs = 17; goto _test_eof

	_test_eof: {}
	if p == eof {
		switch cs {
		case 11:
			goto tr25
		case 5:
			goto tr8
		case 12:
			goto tr27
		case 6:
			goto tr11
		case 13:
			goto tr30
		case 14:
			goto tr27
		case 7:
			goto tr11
		case 15:
			goto tr31
		case 16:
			goto tr32
		}
	}

	_out: {}
	}

//line lang_test.rl:109
	s.SetState(cs, ts, te, act)
	return p, pe
}