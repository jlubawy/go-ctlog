// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/jlubawy/go-ctlog/cmodule"
	"github.com/jlubawy/go-ctlog/ctlog"
)

type TlogInfo struct {
	Date    time.Time      `json:"date"`
	Modules []ctlog.Module `json:"modules"`
}

type ModulesInfo struct {
	Date        time.Time        `json:"date"`
	SearchPaths []string         `json:"searchPaths"`
	Modules     []cmodule.Module `json:"modules"`
}

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
				infof("Must provide a module JSON file.\n")
				fs.Usage()
				os.Exit(1)
			}

			if fs.NArg() != 1 {
				infof("Only accepts one module JSON file.\n")
				fs.Usage()
				os.Exit(1)
			}

			f, err := os.Open(fs.Arg(0))
			if err != nil {
				fatalf("Error opening modules JSON file: %v\n", err)
			}

			var info ModulesInfo
			if err := json.NewDecoder(f).Decode(&info); err != nil {
				f.Close()
				fatalf("Error decoding modules JSON: %v\n", err)
			}
			f.Close()

			var tlogInfo = TlogInfo{
				Date:    time.Now().UTC(),
				Modules: make([]ctlog.Module, 0),
			}

			for _, module := range info.Modules {
				var mf *os.File
				mf, err = os.Open(module.Path)
				if err != nil {
					fatalf("Error opening module file: %v\n", err)
				}

				lines, err := ctlog.FindLines(mf)
				if err != nil {
					mf.Close()
					fatalf("Error finding module lines: %v\n", err)
				}
				mf.Close()

				tlogInfo.Modules = append(tlogInfo.Modules, ctlog.Module{
					Index: module.Index,
					Name:  module.Name,
					Path:  module.Path,
					Lines: lines,
				})
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(&tlogInfo); err != nil {
				fatalf("Error encoding JSON: %v\n", err)
			}
		},
		HelpFn: func() {

		},
	},
	{
		Name: "log",
		CmdFn: func(args []string) {
			fs := flag.NewFlagSet("log", flag.ExitOnError)
			fs.Parse(args)

			if fs.NArg() == 0 {
				infof("Must provide a tokenized logging JSON file.\n")
				fs.Usage()
				os.Exit(1)
			}

			if fs.NArg() != 1 {
				infof("Only accepts one tokenized logging JSON file.\n")
				fs.Usage()
				os.Exit(1)
			}

			f, err := os.Open(fs.Arg(0))
			if err != nil {
				fatalf("Error opening modules JSON file: %v\n", err)
			}

			var tlogInfo TlogInfo
			if err := json.NewDecoder(f).Decode(&tlogInfo); err != nil {
				f.Close()
				fatalf("Error decoding modules JSON: %v\n", err)
			}
			f.Close()

			tx := ctlog.NewTranslator(tlogInfo.Modules)
			s := ctlog.NewScanner(os.Stdin)
			for s.Scan() {
				out, ok, err := ctlog.ParseOutput(s.Bytes())
				if err != nil {
					fatalf("Error parsing tokenized logging output: %v\n", err)
				}
				if ok {
					s, err := tx.Translate(out)
					if err != nil {
						fatalf("Error translating tokenized logging output: %v\n", err)
					}
					fmt.Println(s)
				} else {
					fmt.Println(s.Text())
				}
			}
			if err := s.Err(); err != nil {
				fatalf("Error scanning stding: %v\n", err)
			}
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
	infof(`Usage: ctlog command [options]

Available commands:

    json            walk C source directories and output JSON tokenized logging dictionary
    log             translate raw output using the provided dictionary

Use "ctlog help [command]" for more information about that command.
`)
}

func infof(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func fatalf(format string, args ...interface{}) {
	infof(format, args...)
	os.Exit(1)
}

var templHeader = template.Must(template.New("").Parse(`/**
 * Auto-generated module index definitions for a given project.
 */

// Generated on: {{.Date}}
// Using search paths:
{{range .SearchPaths}}{{printf "//   - %s" .}}
{{end}}

#ifndef MODULE_INDICES_H
#define MODULE_INDICES_H

/*==============================================================================
 *                                   Defines
 *============================================================================*/
/*============================================================================*/
{{range $moduleIdx, $module := .Modules}}#define MODULE_INDEX_{{printf "%-32s" $module.Name}}  ({{$module.Index}})
{{end}}

#endif /* MODULE_INDICES_H */
`))
