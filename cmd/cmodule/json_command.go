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

	"github.com/jlubawy/go-cli"
	"github.com/jlubawy/go-ctlog/cmodule"
)

type JSONOptions struct {
	Compact bool
	Output  string
}

var jsonOptions JSONOptions

var jsonCommand = cli.Command{
	Name:             "json",
	ShortDescription: "walk C source directories and output JSON module info",
	Description:      "Walk C source directories and output JSON module info.",
	ShortUsage:       "[-output output] [directories...]",
	SetupFlags: func(fs *flag.FlagSet) {
		fs.BoolVar(&jsonOptions.Compact, "compact", false, "output compact JSON")
		fs.StringVar(&jsonOptions.Output, "output", "", "output file or stdout if empty")
	},
	Run: func(args []string) {
		if len(args) == 0 {
			cli.Info("Must provide at least one directory.\n")
			os.Exit(1)
		}

		modules, err := cmodule.WalkDirs(args...)
		if err != nil {
			cli.Fatalf("Error walking directories: %v\n", err)
		}

		sps := make([]string, len(args))
		for i := 0; i < len(sps); i++ {
			cp, err := cmodule.PathAbsToSlash(args[i])
			if err != nil {
				cli.Fatalf("Error cleaning search path: %v\n", err)
			}
			sps[i] = cp
		}

		info := ModulesInfo{
			Date:        time.Now().UTC(),
			SearchPaths: sps,
			Modules:     modules,
		}

		var w io.Writer
		if jsonOptions.Output == "" {
			w = os.Stdout
		} else {
			f, err := os.OpenFile(jsonOptions.Output, os.O_CREATE|os.O_WRONLY, 0664)
			if err != nil {
				cli.Fatalf("Error opening output file: %v\n", err)
			}
			defer f.Close()
			w = f
		}

		enc := json.NewEncoder(w)
		if !jsonOptions.Compact {
			enc.SetIndent("", "  ")
		}

		if err := enc.Encode(&info); err != nil {
			cli.Fatalf("Error encoding JSON: %v\n", err)
		}
	},
}
