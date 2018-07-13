
//line next.rl:1
package lexer

import (
    "io"

    "github.com/db47h/monkey/token"
)


//line next.rl:82



//line next.go:17
const monkey_start int = 10
const monkey_error int = 0

const monkey_en_c_comment int = 8
const monkey_en_main int = 10


//line next.rl:85

func (l *Lexer) init() {
    
//line next.go:29
	{
	 l.cs = monkey_start
	 l.ts = 0
	 l.te = 0
	 l.act = 0
	}

//line next.rl:88
}

// Next returns the next token in the input stream.
//
func (l *Lexer) Next() token.Token {
    for l.count == 0 {
        p, pe, eof, err := l.updateBuffer()

        
//line next.go:47
	{
	if p == pe {
		goto _test_eof
	}
	switch  l.cs {
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
//line next.rl:49
 l.te = p+1
{
        l.emit(l.ts, token.String, l.tokenString())
	}
	goto st10
tr6:
//line next.rl:43
 l.te = p+1
{
        l.emit(l.ts, token.Char, l.tokenString())
	}
	goto st10
tr8:
//line next.rl:31
p = ( l.te) - 1
{
		l.emit(l.ts, token.Symbol, l.data[l.ts]);
	}
	goto st10
tr10:
//line next.rl:15
 l.newline(p) 
//line next.rl:59
 l.te = p+1

	goto st10
tr11:
//line next.rl:65
p = ( l.te) - 1
{
        l.emit(l.ts, token.Int, l.tokenString())
	}
	goto st10
tr18:
//line next.rl:54
 l.te = p+1

	goto st10
tr19:
//line next.rl:15
 l.newline(p) 
//line next.rl:54
 l.te = p+1

	goto st10
tr20:
//line next.rl:31
 l.te = p+1
{
		l.emit(l.ts, token.Symbol, l.data[l.ts]);
	}
	goto st10
tr25:
//line next.rl:31
 l.te = p
p--
{
		l.emit(l.ts, token.Symbol, l.data[l.ts]);
	}
	goto st10
tr26:
//line next.rl:61
 l.te = p+1
{ {goto st8 } }
	goto st10
tr27:
//line next.rl:65
 l.te = p
p--
{
        l.emit(l.ts, token.Int, l.tokenString())
	}
	goto st10
tr30:
//line next.rl:71
 l.te = p
p--
{
        l.emit(l.ts, token.Float, l.tokenString())
	}
	goto st10
tr31:
//line next.rl:77
 l.te = p
p--
{
        l.emit(l.ts, token.Int, l.tokenString())
	}
	goto st10
tr32:
//line next.rl:37
 l.te = p
p--
{
        l.emit(l.ts, token.Ident, l.tokenString())
	}
	goto st10
	st10:
//line NONE:1
 l.ts = 0

		if p++; p == pe {
			goto _test_eof10
		}
	st_case_10:
//line NONE:1
 l.ts = p

//line next.go:201
		switch  l.data[p] {
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
		case  l.data[p] < 65:
			switch {
			case  l.data[p] < 49:
				if 33 <=  l.data[p] &&  l.data[p] <= 46 {
					goto tr20
				}
			case  l.data[p] > 57:
				if 58 <=  l.data[p] &&  l.data[p] <= 64 {
					goto tr20
				}
			default:
				goto tr23
			}
		case  l.data[p] > 90:
			switch {
			case  l.data[p] < 97:
				if 91 <=  l.data[p] &&  l.data[p] <= 96 {
					goto tr20
				}
			case  l.data[p] > 122:
				if 123 <=  l.data[p] &&  l.data[p] <= 126 {
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
//line next.rl:15
 l.newline(p) 
	goto st1
	st1:
		if p++; p == pe {
			goto _test_eof1
		}
	st_case_1:
//line next.go:256
		switch  l.data[p] {
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
		if  l.data[p] == 10 {
			goto tr1
		}
		goto st1
tr5:
//line next.rl:15
 l.newline(p) 
	goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line next.go:284
		switch  l.data[p] {
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
		if  l.data[p] == 10 {
			goto tr5
		}
		goto st3
tr21:
//line NONE:1
 l.te = p+1

	goto st11
	st11:
		if p++; p == pe {
			goto _test_eof11
		}
	st_case_11:
//line next.go:313
		switch  l.data[p] {
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
		if  l.data[p] == 10 {
			goto tr10
		}
		goto st5
tr22:
//line NONE:1
 l.te = p+1

	goto st12
	st12:
		if p++; p == pe {
			goto _test_eof12
		}
	st_case_12:
//line next.go:340
		switch  l.data[p] {
		case 46:
			goto st6
		case 120:
			goto st7
		}
		if 48 <=  l.data[p] &&  l.data[p] <= 57 {
			goto tr23
		}
		goto tr27
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
		if 48 <=  l.data[p] &&  l.data[p] <= 57 {
			goto st13
		}
		goto tr11
	st13:
		if p++; p == pe {
			goto _test_eof13
		}
	st_case_13:
		if 48 <=  l.data[p] &&  l.data[p] <= 57 {
			goto st13
		}
		goto tr30
tr23:
//line NONE:1
 l.te = p+1

	goto st14
	st14:
		if p++; p == pe {
			goto _test_eof14
		}
	st_case_14:
//line next.go:379
		if  l.data[p] == 46 {
			goto st6
		}
		if 48 <=  l.data[p] &&  l.data[p] <= 57 {
			goto tr23
		}
		goto tr27
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
		switch {
		case  l.data[p] < 65:
			if 48 <=  l.data[p] &&  l.data[p] <= 57 {
				goto st15
			}
		case  l.data[p] > 70:
			if 97 <=  l.data[p] &&  l.data[p] <= 102 {
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
		case  l.data[p] < 65:
			if 48 <=  l.data[p] &&  l.data[p] <= 57 {
				goto st15
			}
		case  l.data[p] > 70:
			if 97 <=  l.data[p] &&  l.data[p] <= 102 {
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
		if  l.data[p] == 95 {
			goto st16
		}
		switch {
		case  l.data[p] < 65:
			if 48 <=  l.data[p] &&  l.data[p] <= 57 {
				goto st16
			}
		case  l.data[p] > 90:
			if 97 <=  l.data[p] &&  l.data[p] <= 122 {
				goto st16
			}
		default:
			goto st16
		}
		goto tr32
tr15:
//line next.rl:15
 l.newline(p) 
	goto st8
	st8:
//line NONE:1
 l.ts = 0

		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
//line next.go:456
		switch  l.data[p] {
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
		switch  l.data[p] {
		case 10:
			goto tr15
		case 42:
			goto st9
		case 47:
			goto tr17
		}
		goto st8
tr17:
//line next.rl:19
{goto st10 }
	goto st17
	st17:
		if p++; p == pe {
			goto _test_eof17
		}
	st_case_17:
//line next.go:487
		goto st0
st_case_0:
	st0:
		 l.cs = 0
		goto _out
	st_out:
	_test_eof10:  l.cs = 10; goto _test_eof
	_test_eof1:  l.cs = 1; goto _test_eof
	_test_eof2:  l.cs = 2; goto _test_eof
	_test_eof3:  l.cs = 3; goto _test_eof
	_test_eof4:  l.cs = 4; goto _test_eof
	_test_eof11:  l.cs = 11; goto _test_eof
	_test_eof5:  l.cs = 5; goto _test_eof
	_test_eof12:  l.cs = 12; goto _test_eof
	_test_eof6:  l.cs = 6; goto _test_eof
	_test_eof13:  l.cs = 13; goto _test_eof
	_test_eof14:  l.cs = 14; goto _test_eof
	_test_eof7:  l.cs = 7; goto _test_eof
	_test_eof15:  l.cs = 15; goto _test_eof
	_test_eof16:  l.cs = 16; goto _test_eof
	_test_eof8:  l.cs = 8; goto _test_eof
	_test_eof9:  l.cs = 9; goto _test_eof
	_test_eof17:  l.cs = 17; goto _test_eof

	_test_eof: {}
	if p == eof {
		switch  l.cs {
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

//line next.rl:97

        l.shiftBuffer(p, pe)

        switch err {
        case nil:
        case io.EOF:
            l.emit(0, token.EOF, "EOF")
        default:
            l.emit(0, token.Error, err)
        }
    }

    return l.pop()
}