package lexer_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/db47h/monkey/lexer"
	"github.com/db47h/monkey/token"
)

func TestNext(t *testing.T) {
	// input := "a = 4;\nccc := xyz + 17\n"
	// r := strings.NewReader(input)

	r, err := os.Open("/home/denis/devel/ragel/examples/clang.c")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	l := lexer.New(r)
	for {
		tok := l.Next()
		// t.Logf("tok: %v", tok)
		if tok.Type == token.EOF || tok.Type == token.Error {
			break
		}
	}

	tok := l.Next()
	t.Logf("tok: %v", tok)
	t.Logf("Copy: %d", l.Copies)
	t.Logf("Reads: %d", l.Reads)
	t.Logf("QSize: %d", l.QSize())
}

func BenchmarkNext(b *testing.B) {
	input := "a = 4; ccc := xyz + 17\n"
	r := strings.NewReader(input)
	// f, err := os.Open("/home/denis/devel/ragel/examples/clang.c")
	// if err != nil {
	// 	b.Fatal(err)
	// }
	// defer f.Close()
	// data, err := ioutil.ReadAll(f)
	// if err != nil {
	// 	b.Fatal(err)
	// }
	// r := bytes.NewReader(data)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Seek(0, io.SeekStart)
		l := lexer.New(r)
		l.Reset()
		for {
			tok := l.Next()
			if tok.Type == token.EOF || tok.Type == token.Error {
				break
			}
		}
	}
}
