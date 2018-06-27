// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/jlubawy/go-ctlog/cli"

var program = cli.Program{
	Name: "ctlog",
	Commands: []cli.Command{
		dictCommand,
		logCommand,
	},
}

func main() { program.Run() }
