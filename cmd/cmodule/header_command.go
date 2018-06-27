// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/jlubawy/go-cli"
	"github.com/jlubawy/go-ctlog/cmodule"
)

type HeaderOptions struct {
	Output string
}

var headerOptions HeaderOptions

var headerCommand = cli.Command{
	Name:             "header",
	ShortDescription: "create a C header file from the provided cmodules JSON",
	Description:      "Create a C header file from the provided cmodules JSON.",
	ShortUsage:       "[-output output] [cmodule JSON]",
	SetupFlags: func(fs *flag.FlagSet) {
		fs.StringVar(&headerOptions.Output, "output", "", "output file or stdout if empty")
	},
	Run: func(args []string) {
		if len(args) == 0 {
			cli.Info("Must provide a module JSON file.\n")
		}

		if len(args) != 1 {
			cli.Info("Only accepts one module JSON file.\n")
		}

		f, err := os.Open(args[0])
		if err != nil {
			cli.Fatalf("Error opening modules JSON file: %v\n", err)
		}
		defer f.Close()

		var info ModulesInfo
		if err := json.NewDecoder(f).Decode(&info); err != nil {
			cli.Fatalf("Error decoding JSON: %v\n", err)
		}

		var w io.Writer
		if headerOptions.Output == "" {
			w = os.Stdout
		} else {
			f, err := os.OpenFile(headerOptions.Output, os.O_CREATE|os.O_WRONLY, 0664)
			if err != nil {
				cli.Fatalf("Error opening output file: %v\n", err)
			}
			defer f.Close()
			w = f
		}

		if err := templHeader.Execute(w, &info); err != nil {
			cli.Fatalf("Error executing template: %v\n", err)
		}
	},
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

type ModulesInfo struct {
	Date        time.Time        `json:"date"`
	SearchPaths []string         `json:"searchPaths"`
	Modules     []cmodule.Module `json:"modules"`
}
