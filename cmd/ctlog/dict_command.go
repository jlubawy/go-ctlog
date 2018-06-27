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

	"github.com/jlubawy/go-ctlog/cli"
	"github.com/jlubawy/go-ctlog/cmodule"
	"github.com/jlubawy/go-ctlog/ctlog"
)

type DictOptions struct {
	Compact bool
	Output  string
}

var dictOptions DictOptions

var dictCommand = cli.Command{
	Name:             "dict",
	ShortDescription: "create tokenized logging dictionary from a cmodule JSON file",
	Description:      "Dict creates a tokenized logging dictionary from the provided cmodule JSON file.",
	ShortUsage:       "[-output output] [cmodule JSON]",
	SetupFlags: func(fs *flag.FlagSet) {
		fs.BoolVar(&dictOptions.Compact, "compact", false, "output compact JSON")
		fs.StringVar(&dictOptions.Output, "output", "", "output file or stdout if empty")
	},
	Run: func(args []string) {
		if len(args) == 0 {
			cli.Fatal("Must provide a cmodule JSON file.\n")
		}

		if len(args) != 1 {
			cli.Fatal("Only accepts one cmodule JSON file.\n")
		}

		f, err := os.Open(args[0])
		if err != nil {
			cli.Fatalf("Error opening modules JSON file: %v\n", err)
		}

		var info ModulesInfo
		if err := json.NewDecoder(f).Decode(&info); err != nil {
			f.Close()
			cli.Fatalf("Error decoding modules JSON: %v\n", err)
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
				cli.Fatalf("Error opening module file: %v\n", err)
			}

			lines, err := ctlog.FindLines(mf)
			if err != nil {
				mf.Close()
				cli.Fatalf("Error finding module lines: %v\n", err)
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
		if dictOptions.Output == "" {
			w = os.Stdout
		} else {
			f, err := os.OpenFile(dictOptions.Output, os.O_CREATE|os.O_WRONLY, 0664)
			if err != nil {
				cli.Fatalf("Error opening output file: %v\n", err)
			}
			defer f.Close()
			w = f
		}

		enc := json.NewEncoder(w)
		if !dictOptions.Compact {
			enc.SetIndent("", "  ")
		}

		if err := enc.Encode(&tlogInfo); err != nil {
			cli.Fatalf("Error encoding JSON: %v\n", err)
		}
	},
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
