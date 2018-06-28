// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/jlubawy/go-cli"

var program = cli.Program{
	Name:        "cmodule",
	Description: "Cmodule is a program for managing modularized C projects.",
	Commands: []cli.Command{
		headerCommand,
		jsonCommand,
	},
}

func main() { program.RunAndExit() }
