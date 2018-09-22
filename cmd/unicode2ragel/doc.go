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

/*
unicode2ragel generates ragel machines for UTF-8 input based on the unicode
ranges defined in Go's unicode package.


Install

Run the following command to install unicode2ragel:

	go get github.com/db47h/ragel/cmd/unicode2ragel


Usage

Usage of unicode2ragel:

	-cat list
		comma separated list unicode categories to generate (default Letter,Lower,Upper)
	-o filename
		output filename. Writes to stdout if not set

unicode2ragel will generate one ragel state machine per specified category.
The state machine name will be of the form:

	"u" + strings.ToLower(categoryName)

where categoryName is the name as specified on the command line.


Supported unicode character categories

unicode2ragel supports the following unicode character categories:

	Other:  unicode.C,
	Letter: unicode.L,
	Lower:  unicode.Ll,
	Title:  unicode.Lt,
	Upper:  unicode.Lu,
	Mark:   unicode.M,
	Number: unicode.N,
	Digit:  unicode.Nd,
	Punct:  unicode.P,
	Symbol: unicode.S,
	Space:  unicode.Z,

Additionally, all categories from the unicode.Categories map are supported.


Example

In the github.com/db47h/ragel package, the example ragel file lang_test.rl
includes a utf8.rl file that is generated with the following command:

	unicode2ragel -cat Letter -o utf8.rl

Here's an excerpt from the generated file:

	%%{
		machine UTF8;
		uletter =
			0x41..0x5A
			| 0x61..0x7A
			| 0xC2 0xAA
			# … lots of rules removed for brevity …
			| 0xF0 0xAF 0xA0..0xA7 0x80..0xBF
			| 0xF0 0xAF 0xA8 0x80..0x9D
			;
	}%%

The "uletter" state machine can the be used to match utf8 input sequences that
match characters in the unicode L category:

	%%{
		machine lang;
		include UTF8 "utf8.rl";

		# Alpha numberic characters or underscore.
		alnum_u = uletter | digit | '_';

		# Alpha charactres or underscore.
		alpha_u = uletter | '_';
		# …
	}%%
*/
package main
