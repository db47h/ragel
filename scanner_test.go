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
	input := "a = 4;\nçcc = xyz + 17.0\n"
	r := strings.NewReader(input)
	l := ragel.New("", r, fsm{})

	res := []struct {
		pos string
		typ ragel.Token
		lit string
	}{
		{":1:1", Ident, "a"},
		{":1:3", Symbol, "="},
		{":1:5", Int, "4"},
		{":1:6", Symbol, ";"},
		{":2:1", Ident, "çcc"},
		{":2:6", Symbol, "="},
		{":2:8", Ident, "xyz"},
		{":2:12", Symbol, "+"},
		{":2:14", Float, "17.0"},
		{":3:1", ragel.EOF, "EOF"},
	}

	for i := 0; i < len(res); i++ {
		offset, tok, lit := l.Next()
		e := fmt.Sprintf("%s: %v %v", res[i].pos, res[i].typ, res[i].lit)
		g := fmt.Sprintf("%s: %v %v", l.Pos(offset), tok, lit)
		if g != e {
			t.Errorf("Expected %q, got %q", e, g)
		}
	}

	_, tok, _ := l.Next()
	if tok != ragel.EOF {
		t.Errorf("Expected EOF, got %v", tok)
	}
}

func BenchmarkNext(b *testing.B) {
	input := "a = 4; çcc := xyz + 17\n"
	r := strings.NewReader(input)
	l := ragel.New("INPUT", r, fsm{})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Reset()
		for {
			_, tok, _ := l.Next()
			if tok == ragel.EOF || tok == ragel.Error {
				break
			}
		}
	}
}

func BenchmarkNext_largeishFile(b *testing.B) {
	r, err := os.Open("lang_test.go")
	if err != nil {
		b.Fatal(err)
	}
	l := ragel.New("INPUT", r, fsm{})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Reset()
		for {
			_, tok, _ := l.Next()
			if tok == ragel.EOF || tok == ragel.Error {
				break
			}
		}
	}
}
