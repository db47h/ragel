package lexer_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/db47h/monkey/lexer"
)

func TestNext(t *testing.T) {
	input := "a = 4;\nccc = xyz + 17.0\n"
	r := strings.NewReader(input)
	l := lexer.New("", r)

	res := []string{
		":1:1: Ident a",
		":1:3: Symbol =",
		":1:5: Int 4",
		":1:6: Symbol ;",
		":2:1: Ident ccc",
		":2:5: Symbol =",
		":2:7: Ident xyz",
		":2:11: Symbol +",
		":2:13: Float 17.0",
		":3:1: EOF <nil>",
	}

	for i := 0; i < len(res); i++ {
		tok := l.Next()
		g := fmt.Sprintf("%s: %v %v", l.Pos(tok.Offset), tok.Type, tok.Literal)
		if g != res[i] {
			t.Fatalf("Expected %q, got %q", res[i], g)
		}
	}

	tok := l.Next()
	if tok.Type != lexer.EOF {
		t.Fatalf("Expected EOF, got %v", tok)
	}
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
			if tok.Type == lexer.EOF || tok.Type == lexer.Error {
				break
			}
		}
	}
}

func BenchmarkNext_largeishFile(b *testing.B) {
	r, err := os.Open("lang.go")
	if err != nil {
		b.Fatal(err)
	}
	l := lexer.New("INPUT", r)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Reset()
		for {
			tok := l.Next()
			if tok.Type == lexer.EOF || tok.Type == lexer.Error {
				break
			}
		}
	}
}
