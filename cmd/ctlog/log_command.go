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

	"github.com/jlubawy/go-ctlog/cli"
	"github.com/jlubawy/go-ctlog/ctlog"
)

type LogOptions struct {
	Output string
}

var logOptions LogOptions

var logCommand = cli.Command{
	Name:             "log",
	ShortDescription: "translate tokenized logging output using the provided dictionary",
	Description:      "Log translates tokenized logging output using the provided dictionary.",
	ShortUsage:       "[-output output] [dictionary JSON]",
	SetupFlags: func(fs *flag.FlagSet) {
		fs.StringVar(&logOptions.Output, "output", "", "output file or stdout if empty")
	},
	Run: func(args []string) {
		if len(args) == 0 {
			cli.Fatal("Must provide a dictionary JSON file.\n")
		}

		if len(args) != 1 {
			cli.Fatal("Only accepts one dictionary JSON file.\n")
		}

		f, err := os.Open(args[0])
		if err != nil {
			cli.Fatalf("Error opening dictionary JSON file: %v\n", err)
		}

		var tlogInfo TlogInfo
		if err := json.NewDecoder(f).Decode(&tlogInfo); err != nil {
			f.Close()
			cli.Fatalf("Error decoding dictionary JSON: %v\n", err)
		}
		f.Close()

		var w io.Writer
		if logOptions.Output == "" {
			w = os.Stdout
		} else {
			f, err := os.OpenFile(logOptions.Output, os.O_CREATE|os.O_WRONLY, 0664)
			if err != nil {
				cli.Fatalf("Error opening output file: %v\n", err)
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
					cli.Fatalf("Error translating tokenized logging output: %v\n", err)
				}
				fmt.Fprintln(w, s)
			}
		}
		if err := s.Err(); err != nil {
			cli.Fatalf("Error scanning stdin: %v\n", err)
		}
	},
}
