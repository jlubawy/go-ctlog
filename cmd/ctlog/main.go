// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/jlubawy/go-cli"

var program = cli.Program{
	Name:        "ctlog",
	Description: "Ctlog is a program for managing tokenized logging projects.",
	Commands: []cli.Command{
		dictCommand,
		logCommand,
	},
}

func main() { program.RunAndExit() }
