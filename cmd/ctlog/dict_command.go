// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"time"

	"github.com/jlubawy/go-ctlog/cmodule"
	"github.com/jlubawy/go-ctlog/ctlog"
)

const dictUsage = `usage: dict [-output output] [cmodule JSON]
Run 'ctlog help dict' for details.
`

const dictHelp = `usage: ctlog dict [-output output] [cmodule JSON]

Dict creates a tokenized logging dictionary from the provided cmodule JSON file.

Options:

    -compact       output compact JSON dictionary
    -output        file to output the stripped source to, or stdout if empty
`

var dictCommand = Command{
	Name: "dict",
	CmdFn: func(args []string) {
		var (
			flagCompact bool
			flagOutput  string
		)

		fs := flag.NewFlagSet("dict", flag.ExitOnError)
		fs.Usage = func() { info(dictUsage) }
		fs.BoolVar(&flagCompact, "compact", false, "output compact JSON")
		fs.StringVar(&flagOutput, "output", "", "file to output to, stdout if empty")
		fs.Parse(args)

		if fs.NArg() == 0 {
			info("Must provide a cmodule JSON file.\n\n")
			fs.Usage()
			os.Exit(1)
		}

		if fs.NArg() != 1 {
			info("Only accepts one cmodule JSON file.\n\n")
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

		enc := json.NewEncoder(w)

		if !flagCompact {
			enc.SetIndent("", "  ")
		}

		if err := enc.Encode(&tlogInfo); err != nil {
			fatalf("Error encoding JSON: %v\n", err)
		}
	},
	HelpFn: func() { info(dictHelp) },
}

type TlogInfo struct {
	Date    time.Time      `json:"date"`
	Modules []ctlog.Module `json:"modules"`
}

type ModulesInfo struct {
	Date        time.Time        `json:"date"`
	SearchPaths []string         `json:"searchPaths"`
	Modules     []cmodule.Module `json:"modules"`
}
