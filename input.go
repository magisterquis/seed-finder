package main

/*
 * input.go
 * Reads intput line or chunks
 * By J. Stuart McMurray
 * Created 20170712
 * Last Modified 20170712
 */

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/boltdb/bolt"
)

/* getInput returns the strings in args, as well as the lines or
NULL-terminated chunks (if nt is true) from the file named fn. */
func getInput(args []string, fn string, nt bool) ([][]byte, error) {

	/* Convert args to byte slices */
	var o [][]byte
	for _, v := range args {
		o = append(o, []byte(v))
	}

	/* If we have no file, we're done */
	if "" == fn {
		return o, nil
	}

	/* File to read from */
	var r io.Reader
	if "-" == fn {
		r = os.Stdin
	} else {
		/* Try to open the file */
		f, err := os.Open(fn)
		if nil != err {
			return nil, err
		}
		defer f.Close()
		r = f
	}

	/* Read lines or chunks from the file */
	sf := bufio.ScanLines
	if nt {
		sf = ScanChunks
	}

	/* Set up to read nicely */
	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, bolt.MaxKeySize)
	scanner.Split(sf)

	/* Read chunks */
	for scanner.Scan() {
		o = append(o, scanner.Bytes())
	}

	return o, scanner.Err()
}

/* The following function taken from
https://github.com/golang/go/blob/a1110c39301b21471c27dad0e50cdbe499587fc8/src/bufio/scan.go
with minor modifications.  It is under the following license:

Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
/* ScanChunks is a bufio.SplitFunc that returns \0-terminated chunks */
func ScanChunks(
	data []byte,
	atEOF bool,
) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
