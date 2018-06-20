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
	"github.com/jlubawy/go-ctlog/internal"
)

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
		Name: "header",
		CmdFn: func(args []string) {
			fs := flag.NewFlagSet("header", flag.ExitOnError)
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
			defer f.Close()

			var info ModulesInfo
			if err := json.NewDecoder(f).Decode(&info); err != nil {
				fatalf("Error decoding JSON: %v\n", err)
			}

			if err := templHeader.Execute(os.Stdout, &info); err != nil {
				fatalf("Error executing template: %v\n", err)
			}
		},
		HelpFn: func() {

		},
	},
	{
		Name: "json",
		CmdFn: func(args []string) {
			fs := flag.NewFlagSet("json", flag.ExitOnError)
			fs.Parse(args)

			if fs.NArg() == 0 {
				infof("Must provide at least one directory.\n")
				fs.Usage()
				os.Exit(1)
			}

			modules, err := cmodule.WalkDirs(fs.Args()...)
			if err != nil {
				fatalf("Error walking directories: %v\n", err)
			}

			sps := make([]string, fs.NArg())
			for i := 0; i < len(sps); i++ {
				cp, err := internal.PathAbsToSlash(fs.Arg(i))
				if err != nil {
					fatalf("Error cleaning search path: %v\n", err)
				}
				sps[i] = cp
			}

			info := ModulesInfo{
				Date:        time.Now().UTC(),
				SearchPaths: sps,
				Modules:     modules,
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(&info); err != nil {
				fatalf("Error encoding JSON: %v\n", err)
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
	infof(`Usage: cmodule command [options]

Available commands:

    json            walk C source directories and output JSON module info

Use "cmodule help [command]" for more information about that command.
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

#ifndef CMODULE_INDICES_H
#define CMODULE_INDICES_H

/*==============================================================================
 *                                   Defines
 *============================================================================*/
/*============================================================================*/
{{range $moduleIdx, $module := .Modules}}#define CMODULE_INDEX_{{printf "%-32s" $module.Name}}  ({{$module.Index}})
{{end}}

#endif /* CMODULE_INDICES_H */
`))
