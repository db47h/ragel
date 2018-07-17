
//line csv.rl:1
// Example CSV scanner with comma separated fields.

package csv

import "github.com/db47h/ragel"


//line csv.rl:43



//line csv.rl:47

//line csv.go:17
const csv_start int = 1
const csv_error int = -1

const csv_en_main int = 1


//line csv.rl:48

// here we use a private implemetation for ragel.FSM since the scanner and
// parser are in the same package.
//
type fsm struct {}

func (fsm) Init(s *ragel.Scanner) {
	var cs, ts, te, act int
	
//line csv.go:34
	{
	cs = csv_start
	ts = 0
	te = 0
	act = 0
	}

//line csv.rl:57
	s.SetState(cs, ts, te, act)
}

func (fsm) States() (start, err int) {
	return 1, -1
}

func (fsm) Run(s *ragel.Scanner, p, pe, eof int) (int, int) {
	cs, ts, te, act, data := s.GetState()
	
//line csv.go:53
	{
	if p == pe {
		goto _test_eof
	}
	switch cs {
	case 1:
		goto st_case_1
	case 2:
		goto st_case_2
	case 3:
		goto st_case_3
	case 4:
		goto st_case_4
	case 0:
		goto st_case_0
	case 5:
		goto st_case_5
	case 6:
		goto st_case_6
	case 7:
		goto st_case_7
	case 8:
		goto st_case_8
	case 9:
		goto st_case_9
	case 10:
		goto st_case_10
	}
	goto st_out
tr0:
//line NONE:1
	switch act {
	case 2:
	{p = (te) - 1

        	s.Emit(ts, QuotedString, string(data[ts:te]))
	}
	case 4:
	{p = (te) - 1

        	s.Emit(ts, Number, string(data[ts:te]))
	}
	case 6:
	{p = (te) - 1
 s.Emit(ts, String, string(data[ts:te])) }
	}
	
	goto st1
tr6:
//line csv.rl:10
 s.Newline(p) 
//line csv.rl:26
te = p+1
{
        s.Emit(ts, EOL, "\n")
     }
	goto st1
tr9:
//line csv.rl:15
te = p+1
{
		s.Emit(ts, Comma, string(data[ts:te]));
	}
	goto st1
tr11:
//line csv.rl:40
te = p
p--
{ s.Emit(ts, String, string(data[ts:te])) }
	goto st1
tr12:
//line csv.rl:37
te = p
p--

	goto st1
tr14:
//line csv.rl:21
te = p
p--
{
        	s.Emit(ts, QuotedString, string(data[ts:te]))
	}
	goto st1
tr15:
//line csv.rl:32
te = p
p--
{
        	s.Emit(ts, Number, string(data[ts:te]))
	}
	goto st1
	st1:
//line NONE:1
ts = 0

		if p++; p == pe {
			goto _test_eof1
		}
	st_case_1:
//line NONE:1
ts = p

//line csv.go:157
		switch data[p] {
		case 10:
			goto tr6
		case 32:
			goto st3
		case 34:
			goto tr7
		case 44:
			goto tr9
		}
		switch {
		case data[p] < 43:
			if 9 <= data[p] && data[p] <= 13 {
				goto st3
			}
		case data[p] > 45:
			if 48 <= data[p] && data[p] <= 57 {
				goto st8
			}
		default:
			goto st7
		}
		goto st2
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
		switch data[p] {
		case 32:
			goto tr11
		case 44:
			goto tr11
		}
		if 9 <= data[p] && data[p] <= 13 {
			goto tr11
		}
		goto st2
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
		switch data[p] {
		case 9:
			goto st3
		case 32:
			goto st3
		}
		if 11 <= data[p] && data[p] <= 13 {
			goto st3
		}
		goto tr12
tr7:
//line NONE:1
te = p+1

//line csv.rl:40
act = 6;
	goto st4
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
//line csv.go:223
		switch data[p] {
		case 10:
			goto tr2
		case 32:
			goto st0
		case 34:
			goto st6
		case 44:
			goto st0
		}
		if 9 <= data[p] && data[p] <= 13 {
			goto st0
		}
		goto tr7
tr2:
//line csv.rl:10
 s.Newline(p) 
	goto st0
	st0:
		if p++; p == pe {
			goto _test_eof0
		}
	st_case_0:
//line csv.go:247
		switch data[p] {
		case 10:
			goto tr2
		case 34:
			goto tr3
		}
		goto st0
tr3:
//line NONE:1
te = p+1

//line csv.rl:21
act = 2;
	goto st5
	st5:
		if p++; p == pe {
			goto _test_eof5
		}
	st_case_5:
//line csv.go:267
		if data[p] == 34 {
			goto st0
		}
		goto tr14
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
		switch data[p] {
		case 32:
			goto tr14
		case 34:
			goto tr7
		case 44:
			goto tr14
		}
		if 9 <= data[p] && data[p] <= 13 {
			goto tr14
		}
		goto st2
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
		switch data[p] {
		case 32:
			goto tr11
		case 44:
			goto tr11
		}
		switch {
		case data[p] > 13:
			if 48 <= data[p] && data[p] <= 57 {
				goto st8
			}
		case data[p] >= 9:
			goto tr11
		}
		goto st2
	st8:
		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
		switch data[p] {
		case 32:
			goto tr15
		case 44:
			goto tr15
		case 46:
			goto st7
		case 101:
			goto st9
		}
		switch {
		case data[p] > 13:
			if 48 <= data[p] && data[p] <= 57 {
				goto st8
			}
		case data[p] >= 9:
			goto tr15
		}
		goto st2
	st9:
		if p++; p == pe {
			goto _test_eof9
		}
	st_case_9:
		switch data[p] {
		case 32:
			goto tr11
		case 44:
			goto tr11
		}
		switch {
		case data[p] < 43:
			if 9 <= data[p] && data[p] <= 13 {
				goto tr11
			}
		case data[p] > 45:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr18
			}
		default:
			goto tr17
		}
		goto st2
tr17:
//line NONE:1
te = p+1

//line csv.rl:40
act = 6;
	goto st10
tr18:
//line NONE:1
te = p+1

//line csv.rl:32
act = 4;
	goto st10
	st10:
		if p++; p == pe {
			goto _test_eof10
		}
	st_case_10:
//line csv.go:376
		switch data[p] {
		case 32:
			goto tr0
		case 44:
			goto tr0
		}
		switch {
		case data[p] > 13:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr18
			}
		case data[p] >= 9:
			goto tr0
		}
		goto st2
	st_out:
	_test_eof1: cs = 1; goto _test_eof
	_test_eof2: cs = 2; goto _test_eof
	_test_eof3: cs = 3; goto _test_eof
	_test_eof4: cs = 4; goto _test_eof
	_test_eof0: cs = 0; goto _test_eof
	_test_eof5: cs = 5; goto _test_eof
	_test_eof6: cs = 6; goto _test_eof
	_test_eof7: cs = 7; goto _test_eof
	_test_eof8: cs = 8; goto _test_eof
	_test_eof9: cs = 9; goto _test_eof
	_test_eof10: cs = 10; goto _test_eof

	_test_eof: {}
	if p == eof {
		switch cs {
		case 2:
			goto tr11
		case 3:
			goto tr12
		case 4:
			goto tr11
		case 0:
			goto tr0
		case 5:
			goto tr14
		case 6:
			goto tr14
		case 7:
			goto tr11
		case 8:
			goto tr15
		case 9:
			goto tr11
		case 10:
			goto tr0
		}
	}

	}

//line csv.rl:67
	s.SetState(cs, ts, te, act)
	return p, pe
}
