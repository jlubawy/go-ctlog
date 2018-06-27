// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jlubawy/go-ctlog/ctlog"
)

const logUsage = `usage: log [-output output] [dictionary JSON]
Run 'ctlog help log' for details.
`

const logHelp = `usage: ctlog log [-output output] [dictionary JSON]

Log translates tokenized logging output using the provided dictionary.

Options:

    -output        file to output the translated logging to, or stdout if empty
`

var logCommand = Command{
	Name: "log",
	CmdFn: func(args []string) {
		var flagOutput string

		fs := flag.NewFlagSet("log", flag.ExitOnError)
		fs.Usage = func() { info(logUsage) }
		fs.StringVar(&flagOutput, "output", "", "file to output to, stdout if empty")
		fs.Parse(args)

		if fs.NArg() == 0 {
			info("Must provide a dictionary JSON file.\n\n")
			fs.Usage()
			os.Exit(1)
		}

		if fs.NArg() != 1 {
			info("Only accepts one dictionary JSON file.\n\n")
			fs.Usage()
			os.Exit(1)
		}

		f, err := os.Open(fs.Arg(0))
		if err != nil {
			fatalf("Error opening dictionary JSON file: %v\n", err)
		}

		var tlogInfo TlogInfo
		if err := json.NewDecoder(f).Decode(&tlogInfo); err != nil {
			f.Close()
			fatalf("Error decoding dictionary JSON: %v\n", err)
		}
		f.Close()

		var w io.Writer
		if flagOutput == "" {
			w = os.Stdout
		} else {
			f, err := os.OpenFile(flagOutput, os.O_CREATE|os.O_WRONLY, 0664)
			if err != nil {
				fatalf("Error opening output file: %v\n", err)
			}
			defer f.Close()
			w = f
		}

		tx := ctlog.NewTranslator(tlogInfo.Modules)
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			var out ctlog.Output
			if err := json.Unmarshal(s.Bytes(), &out); err != nil {
				fmt.Fprintln(w, s.Text())
			} else {
				s, err := tx.Translate(&out)
				if err != nil {
					fatalf("Error translating tokenized logging output: %v\n", err)
				}
				fmt.Fprintln(w, s)
			}
		}
		if err := s.Err(); err != nil {
			fatalf("Error scanning stdin: %v\n", err)
		}
	},
	HelpFn: func() { info(logHelp) },
}
