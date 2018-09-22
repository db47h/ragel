// Copyright 2018 Denis Bernard <db047h@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var uCategories = map[string]*unicode.RangeTable{
	"Other":  unicode.C,
	"Letter": unicode.L,
	"Lower":  unicode.Ll,
	"Title":  unicode.Lt,
	"Upper":  unicode.Lu,
	"Mark":   unicode.M,
	"Number": unicode.N,
	"Digit":  unicode.Nd,
	"Punct":  unicode.P,
	"Symbol": unicode.S,
	"Space":  unicode.Z,
}

type categories struct {
	names []string
	rts   []*unicode.RangeTable
}

func (cs *categories) String() string {
	return strings.Join(cs.names, ",")
}

func (cs *categories) Set(s string) error {
	cs.names = strings.Split(strings.Replace(s, " ", "", -1), ",")
	cs.rts = nil
	for _, n := range cs.names {
		if rt := uCategories[n]; rt != nil {
			cs.rts = append(cs.rts, rt)
			continue
		}
		if rt := unicode.Categories[n]; rt != nil {
			cs.rts = append(cs.rts, rt)
			continue
		}
		return fmt.Errorf("unknown category %q", n)
	}
	return nil
}

func main() {
	var (
		outFile    string
		categories = categories{
			names: []string{"Letter", "Lower", "Upper"},
			rts:   []*unicode.RangeTable{unicode.Letter, unicode.Upper, unicode.Lower},
		}
		machineName string
	)

	flag.StringVar(&machineName, "m", "UTF8", "machine `name` to use in the ragel machine and import statements")
	flag.Var(&categories, "cat", "comma separated `list` of unicode categories to generate")
	flag.StringVar(&outFile, "o", "", "output `filename`. Writes to stdout if not set")

	flag.Parse()

	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()

	if outFile != "" && outFile != "-" {
		f, err := os.Create(outFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
		defer f.Close()
		os.Stdout = f
	}

	fmt.Print(`// Code generated by unicode2ragel. DO NOT EDIT.
//
// This file defines `)

	for i, n := range categories.names {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%q", "u"+strings.ToLower(n))
	}

	fmt.Print(`
//
// Generated from the Go unicode package (unicode version `, unicode.Version, `). 
// Make sure that alphtype is set to unsigned char for C, and byte for Go.

%%{
	machine `, machineName, ";\n")

	for i, n := range categories.names {
		newClass("u"+strings.ToLower(n), categories.rts[i]).print()
	}

	fmt.Print("}%%\n")
}

type class struct {
	name   string
	rt     *unicode.RangeTable
	ranges []br
}

func newClass(name string, rt *unicode.RangeTable) *class {
	c := &class{name: name, rt: rt}
	for _, r16 := range rt.R16 {
		c.addRange(rune(r16.Lo), rune(r16.Hi), uint32(r16.Stride))
	}
	for _, r32 := range rt.R32 {
		c.addRange(rune(r32.Lo), rune(r32.Hi), r32.Stride)
	}
	return c
}

func (c *class) print() {
	fmt.Printf("\t%s =\n", c.name)
	for i, cur := range c.ranges {
		if i > 0 {
			fmt.Print("\t\t| ")
		} else {
			fmt.Print("\t\t  ")
		}
		fmt.Println(cur.String())
	}
	fmt.Println("\t\t;")
}

func (c *class) add(p []byte) {
	r := make(br, len(p))
	for i, b := range p {
		r[i].s, r[i].e = b, b
	}
	c.ranges = append(c.ranges, r)
	c.merge()
}

func (c *class) merge() {
	l := len(c.ranges)
	if l < 2 {
		return
	}
	prev := c.ranges[l-2]
	last := c.ranges[l-1]
	if len(prev) != len(last) {
		return
	}

	m := -1 // merge position. All other positions must be =
	for i := 0; i < len(prev); i++ {
		if prev[i] == last[i] {
			continue
		}
		// accept one and only one mergeable position
		if m >= 0 {
			return
		}
		if prev[i].e < 255 && prev[i].e+1 == last[i].s {
			m = i
		}
	}
	if m >= 0 {
		prev[m].e = last[m].e
		c.ranges = c.ranges[:l-1]
		c.merge()
	}
}

type br []struct {
	s, e byte
}

func (b br) String() string {
	var sb strings.Builder
	for i, r := range b {
		if i > 0 {
			sb.WriteRune(' ')
		}
		if r.s == r.e {
			sb.WriteString(fmt.Sprintf("0x%02X", r.s))
		} else {
			sb.WriteString(fmt.Sprintf("0x%02X..0x%02X", r.s, r.e))
		}
	}
	return sb.String()
}

func (c *class) addRange(lo, hi rune, stride uint32) {
	var b = make([]byte, 4)

	inc := int32(stride)
	for r := lo; r <= hi; r += inc {
		n := utf8.EncodeRune(b, r)
		if n == 0 {
			continue
		}
		c.add(b[:n])
	}
}
