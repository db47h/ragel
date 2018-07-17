// Copyright 2018 Denis Bernard <db047h@gmail.com>
// Licensed under the MIT license. See license text in the LICENSE file.

//go:generate ragel -Z -G2 csv.rl

// Package csv implements a simple CSV parser built on top of a ragel scanner.
//
package csv

import (
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/db47h/ragel"
)

// CSV tokens
//
const (
	Number ragel.Token = iota
	String
	QuotedString
	Comma
	EOL
)

// Parser is a simple CSV parser.
//
type Parser struct {
	s *ragel.Scanner
}

// New returns a new parser for the given fileName and io.Reader.
//
func New(fileName string, r io.Reader) *Parser {
	p := &Parser{
		// fsm is defined in csv.go, generated from csv.rl by
		// the command `ragel -Z -G2 csv.rl`
		s: ragel.New(fileName, r, fsm{}),
	}
	return p
}

// NextRecord returns the next record ion the CSV file. At the end
// of the file, it returns io.EOF.
//
func (p *Parser) NextRecord() ([]interface{}, error) {
	var (
		rec []interface{}
		err error
		v   interface{}
	)

	pos, t, lit := p.s.Next()
	if t == ragel.EOF {
		return nil, io.EOF
	}
	if t == EOL {
		return nil, nil
	}

	for {
		switch t {
		case ragel.Error:
			if err == nil {
				err = fmt.Errorf("parse error: %v: %s", p.s.Pos(pos), lit)
			}
		case ragel.EOF, EOL:
			if err != nil {
				return nil, err
			}
			rec = append(rec, v)
			return rec, err
		case Number:
			f := new(big.Float)
			f, _, err = f.Parse(lit, 10)
			v = f
		case String:
			// TODO: we might want to validate any quotes
			v = lit
		case QuotedString:
			// remove quotes
			v = strings.Replace(lit[1:len(lit)-1], "\"\"", "\"", -1)
		case Comma:
			rec = append(rec, v)
			v = nil
		}
		pos, t, lit = p.s.Next()
	}
}