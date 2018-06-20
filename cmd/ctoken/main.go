// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jlubawy/go-ctlog/ctoken"
)

type Command struct {
	Name   string
	CmdFn  func(args []string)
	HelpFn func()
}

var commands = []Command{
	{
		Name: "json",
		CmdFn: func(args []string) {
			fs := flag.NewFlagSet("json", flag.ExitOnError)
			fs.Parse(args)

			if fs.NArg() == 0 {
				infof("Must provide a file to parse.\n")
				fs.Usage()
				os.Exit(1)
			}

			if fs.NArg() > 1 {
				infof("Only one file allowed.\n")
				fs.Usage()
				os.Exit(1)
			}

			f, err := os.Open(fs.Arg(0))
			if err != nil {
				fatalf("Error opening file: %v\n", err)
			}
			defer f.Close()

			tokens := make([]ctoken.Token, 0)
			z := ctoken.NewTokenizer(f)
			for {
				tt := z.Next()

				switch tt {
				case ctoken.TokenTypeError:
					err := z.Err()
					if err == io.EOF {
						goto DONE
					}
					fatalf("Error tokenizing file: %v\n", err)

				case ctoken.TokenTypeComment:
					tokens = append(tokens, z.Comment())

				case ctoken.TokenTypeText:
					tokens = append(tokens, z.Text())
				}
			}

		DONE:
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(&tokens); err != nil {
				fatalf("Error encoding JSON: %v\n", err)
			}
		},
		HelpFn: func() {

		},
	},
	{
		Name: "recreate",
		CmdFn: func(args []string) {
			fs := flag.NewFlagSet("recreate", flag.ExitOnError)
			fs.Parse(args)

			if fs.NArg() == 0 {
				infof("Must provide a file to parse.\n")
				fs.Usage()
				os.Exit(1)
			}

			if fs.NArg() > 1 {
				infof("Only one file allowed.\n")
				fs.Usage()
				os.Exit(1)
			}

			f, err := os.Open(fs.Arg(0))
			if err != nil {
				fatalf("Error opening file: %v\n", err)
			}
			defer f.Close()

			z := ctoken.NewTokenizer(f)
			for {
				tt := z.Next()

				switch tt {
				case ctoken.TokenTypeError:
					err := z.Err()
					if err == io.EOF {
						goto DONE
					}
					fatalf("Error tokenizing file: %v\n", err)

				case ctoken.TokenTypeComment:
					fmt.Fprint(os.Stdout, z.Comment().Data)

				case ctoken.TokenTypeText:
					fmt.Fprint(os.Stdout, z.Text().Data)
				}
			}

		DONE:
		},
		HelpFn: func() {

		},
	},
	{
		Name: "stripcomments",
		CmdFn: func(args []string) {
			fs := flag.NewFlagSet("stripcomments", flag.ExitOnError)
			fs.Parse(args)

			if fs.NArg() == 0 {
				infof("Must provide a file to parse.\n")
				fs.Usage()
				os.Exit(1)
			}

			if fs.NArg() > 1 {
				infof("Only one file allowed.\n")
				fs.Usage()
				os.Exit(1)
			}

			f, err := os.Open(fs.Arg(0))
			if err != nil {
				fatalf("Error opening file: %v\n", err)
			}
			defer f.Close()

			z := ctoken.NewTokenizer(f)
			for {
				tt := z.Next()

				switch tt {
				case ctoken.TokenTypeError:
					err := z.Err()
					if err == io.EOF {
						goto DONE
					}
					fatalf("Error tokenizing file: %v\n", err)

				case ctoken.TokenTypeComment:
					// skip comments

				case ctoken.TokenTypeText:
					fmt.Fprintln(os.Stdout, z.Text())
				}
			}

		DONE:
		},
		HelpFn: func() {

		},
	},
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, cmd := range commands {
		if flag.Arg(0) == "help" {
			flag.Usage()
			os.Exit(1)
		} else if flag.Arg(0) == cmd.Name {
			if flag.NArg() >= 2 {
				if flag.Arg(1) == "help" {
					cmd.HelpFn()
					os.Exit(1)
				}
			}
			cmd.CmdFn(flag.Args()[1:])
			os.Exit(0)
		}
	}

	fatalf("Unknown command \"%s\"\n", flag.Arg(0))
}

func usage() {
	infof(`Usage: ctoken command [options]

Available commands:

	json            log C tokens to a JSON file, can be used for testing
	recreate        tokenize the input file and output it, attempting to recreate
	stripcomments   strip a C source file of comments

Use "ctoken help [command]" for more information about that command.
`)
}

func infof(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func fatalf(format string, args ...interface{}) {
	infof(format, args...)
	os.Exit(1)
}
