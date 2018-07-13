package lexer_test

import (
	"io"
	"strings"
	"testing"

	"github.com/db47h/monkey/lexer"
	"github.com/db47h/monkey/token"
)

func TestNext(t *testing.T) {
	input := `a = 4; ccc := xyz + 17
`
	l := lexer.New(strings.NewReader(input))
	for {
		tok := l.Next()
		t.Logf("tok: %v", tok)
		if tok.Type == token.EOF || tok.Type == token.Error {
			break
		}
	}
}

func BenchmarkNext(b *testing.B) {
	input := "a = 4; ccc := xyz + 17\n"

	r := strings.NewReader(input)

	for i := 0; i < b.N; i++ {
		r.Seek(0, io.SeekStart)
		l := lexer.New(r)
		for {
			tok := l.Next()
			if tok.Type == token.EOF || tok.Type == token.Error {
				break
			}
		}
	}
}
