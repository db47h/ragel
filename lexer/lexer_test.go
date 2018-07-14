package lexer_test

import (
	"strings"
	"testing"

	"github.com/db47h/monkey/lexer"
	"github.com/db47h/monkey/token"
)

func TestNext(t *testing.T) {
	input := "a = 4;\nccc := xyz + 17\n"
	r := strings.NewReader(input)

	l := lexer.New("INPUT", r)
	for {
		tok := l.Next()
		t.Logf("%s: %v", l.Pos(tok.Offset), tok)
		if tok.Type == token.EOF || tok.Type == token.Error {
			break
		}
	}

	tok := l.Next()
	t.Logf("tok: %v", tok)
}

func BenchmarkNext(b *testing.B) {
	input := "a = 4; ccc := xyz + 17\n"
	r := strings.NewReader(input)
	l := lexer.New("INPUT", r)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Reset()
		for {
			tok := l.Next()
			if tok.Type == token.EOF || tok.Type == token.Error {
				break
			}
		}
	}
}
