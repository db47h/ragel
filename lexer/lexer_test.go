package lexer_test

import (
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
