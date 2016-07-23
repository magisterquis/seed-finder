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
	"log"
	"os"

	"github.com/boltdb/bolt"
)

/* getInput returns a channel from which can be read the strings in args as
well as the lines or null-terminated chunks (if nt is true) in the file named
fn. */
func getInput(args []string, fn string, nt bool) (<-chan []byte, error) {
	/* Convert args to byte slices */
	var o [][]byte
	for _, v := range args {
		o = append(o, []byte(v))
	}
	var r *os.File
	/* Open the file if we have one */
	if "" != fn {
		/* File to read from */
		if "-" == fn {
			r = os.Stdin
		} else {
			/* Try to open the file */
			var err error
			r, err = os.Open(fn)
			if nil != err {
				return nil, err
			}
		}
	}
	/* Channel to send strings on */
	ch := make(chan []byte)
	/* Read into the channel */
	go inputToChannel(o, r, nt, ch)
	return ch, nil

}

/* inputToChannel puts the slices from args on a channel, as well as the
lines or chunks from r */
func inputToChannel(args [][]byte, r *os.File, nt bool, ch chan<- []byte) {
	if nil != r {
		defer r.Close()
	}
	defer close(ch)

	/* Send args on the channel */
	for _, a := range args {
		ch <- a
	}

	/* Don't bother if we didn't have a file */
	if nil == r {
		return
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
		ch <- scanner.Bytes()
	}
	if err := scanner.Err(); nil != err {
		log.Printf("Unable to read input: %v", err)
	}
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
