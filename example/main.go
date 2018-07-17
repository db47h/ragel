// Copyright 2018 Denis Bernard <db047h@gmail.com>
// Licensed under the MIT license. See license text in the LICENSE file.

// This example shows how to build a parser on top of a ragel scanner
// to parse CSV files.
//
package main

import (
	"io"
	"log"
	"strings"

	"github.com/db47h/ragel/example/csv"
)

func main() {
	data := `12, -17.5, 1e+7, 1.5e-8
			 unquoted_id, another, $cost, chg%
			 "a", "", """", "a ""quote""
with newline"`
	r := strings.NewReader(data)
	p := csv.New("", r)

	for {
		rec, err := p.NextRecord()
		log.Printf("%#v", rec)
		if err == io.EOF {
			break
		}
	}
}
