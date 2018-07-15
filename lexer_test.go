//go:generate ragel -Z -G2 lang_test.rl

package ragel_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/db47h/ragel"
)

func TestNext(t *testing.T) {
	input := "a = 4;\nccc = xyz + 17.0\n"
	r := strings.NewReader(input)
	l := ragel.New("", r, fsm{})

	res := []struct {
		pos string
		typ ragel.Type
		lit string
	}{
		{":1:1", Ident, "a"},
		{":1:3", Symbol, "="},
		{":1:5", Int, "4"},
		{":1:6", Symbol, ";"},
		{":2:1", Ident, "ccc"},
		{":2:5", Symbol, "="},
		{":2:7", Ident, "xyz"},
		{":2:11", Symbol, "+"},
		{":2:13", Float, "17.0"},
		{":3:1", ragel.EOF, ""},
	}

	for i := 0; i < len(res); i++ {
		tok := l.Next()
		e := fmt.Sprintf("%s: %v %v", res[i].pos, res[i].typ, []byte(res[i].lit))
		if tok.Literal == nil {
			tok.Literal = []byte(nil)
		}
		g := fmt.Sprintf("%s: %v %v", l.Pos(tok.Offset), tok.Type, tok.Literal)
		if g != e {
			t.Fatalf("Expected %q, got %q", e, g)
		}
	}

	tok := l.Next()
	if tok.Type != ragel.EOF {
		t.Fatalf("Expected EOF, got %v", tok)
	}
}

func BenchmarkNext(b *testing.B) {
	input := "a = 4; ccc := xyz + 17\n"
	r := strings.NewReader(input)
	l := ragel.New("INPUT", r, fsm{})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Reset()
		for {
			tok := l.Next()
			if tok.Type == ragel.EOF || tok.Type == ragel.Error {
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
	l := ragel.New("INPUT", r, fsm{})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Reset()
		for {
			tok := l.Next()
			if tok.Type == ragel.EOF || tok.Type == ragel.Error {
				break
			}
		}
	}
}
