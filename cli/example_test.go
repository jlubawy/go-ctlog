// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jlubawy/go-ctlog/cli"
)

type EchoOptions struct {
	Prefix string
}

var echoOptions EchoOptions

var program = cli.Program{
	Name: "example",
	Commands: []cli.Command{
		{
			Name:             "echo",
			ShortDescription: "echo the provided text",
			Description:      "Echo the provided text.",
			ShortUsage:       "[text to echo]",
			SetupFlags: func(fs *flag.FlagSet) {
				fs.StringVar(&echoOptions.Prefix, "prefix", "", "prefix to be added to the text")
			},
			Run: func(args []string) {
				var text string
				if len(args) >= 1 {
					text = strings.Join(args, " ")
				}
				if echoOptions.Prefix == "" {
					fmt.Println(text)
				} else {
					fmt.Printf("%s: %s\n", echoOptions.Prefix, text)
				}
			},
		},
	},
}

func Example_noargs() {
	cli.Writer = os.Stdout
	args := []string{
		"example",
	}
	program.Run(args)

	// Output:
	// Usage: example command [options]
	//
	// Available commands:
	//
	//     echo       echo the provided text
	//
	// Use "example help [command]" for more Information about that command.
}

func Example_help() {
	cli.Writer = os.Stdout
	args := []string{
		"example",
		"help",
		"echo",
	}
	program.Run(args)

	// Output:
	// Usage: example echo [text to echo]
	//
	// Echo the provided text.
	//
	// Options:
	//
	//     -prefix     prefix to be added to the text (default=)
}

func Example_run() {
	cli.Writer = os.Stdout
	args := []string{
		"example",
		"echo",
		"-prefix=Some prefix",
		"Some random text",
	}
	program.Run(args)

	// Output:
	// Some prefix: Some random text
}
