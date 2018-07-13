
//line next.rl:1
package lexer

import (
    "errors"
    "io"

    "github.com/db47h/monkey/token"
)


//line next.rl:83



//line next.go:18
const monkey_start int = 10
const monkey_error int = 0

const monkey_en_c_comment int = 8
const monkey_en_main int = 10


//line next.rl:86

func (l *Lexer) init() {
    go l.run()
}


func (l *Lexer) run()  {
    var cs, act, ts, te int

    
//line next.go:37
	{
	cs = monkey_start
	ts = 0
	te = 0
	act = 0
	}

//line next.rl:96

    _ = act

    var have int
    var data [10]byte

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

        
//line next.go:82
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
//line next.rl:50
te = p+1
{
        l.emit(ts, token.String, string(data[ts:te]))
	}
	goto st10
tr6:
//line next.rl:44
te = p+1
{
        l.emit(ts, token.Char, string(data[ts:te]))
	}
	goto st10
tr8:
//line next.rl:32
p = (te) - 1
{
		l.emit(ts, token.Symbol, data[ts]);
	}
	goto st10
tr10:
//line next.rl:16
 l.newline(p) 
//line next.rl:60
te = p+1

	goto st10
tr11:
//line next.rl:66
p = (te) - 1
{
        l.emit(ts, token.Int, string(data[ts:te]))
	}
	goto st10
tr18:
//line next.rl:55
te = p+1

	goto st10
tr19:
//line next.rl:16
 l.newline(p) 
//line next.rl:55
te = p+1

	goto st10
tr20:
//line next.rl:32
te = p+1
{
		l.emit(ts, token.Symbol, data[ts]);
	}
	goto st10
tr25:
//line next.rl:32
te = p
p--
{
		l.emit(ts, token.Symbol, data[ts]);
	}
	goto st10
tr26:
//line next.rl:62
te = p+1
{ {goto st8 } }
	goto st10
tr27:
//line next.rl:66
te = p
p--
{
        l.emit(ts, token.Int, string(data[ts:te]))
	}
	goto st10
tr30:
//line next.rl:72
te = p
p--
{
        l.emit(ts, token.Float, string(data[ts:te]))
	}
	goto st10
tr31:
//line next.rl:78
te = p
p--
{
        l.emit(ts, token.Int, string(data[ts:te]))
	}
	goto st10
tr32:
//line next.rl:38
te = p
p--
{
        l.emit(ts, token.Ident, string(data[ts:te]))
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

//line next.go:236
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
//line next.rl:16
 l.newline(p) 
	goto st1
	st1:
		if p++; p == pe {
			goto _test_eof1
		}
	st_case_1:
//line next.go:291
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
//line next.rl:16
 l.newline(p) 
	goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line next.go:319
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
//line next.go:348
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
//line next.go:375
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
//line next.go:414
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
//line next.rl:16
 l.newline(p) 
	goto st8
	st8:
//line NONE:1
ts = 0

		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
//line next.go:491
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
//line next.rl:20
{goto st10 }
	goto st17
	st17:
		if p++; p == pe {
			goto _test_eof17
		}
	st_case_17:
//line next.go:522
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

//line next.rl:132

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