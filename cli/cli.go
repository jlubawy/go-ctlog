// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package cli provides a simple way of creating new command-line interface (CLI)
programs. It is built using nothing but standard packages and is based on the
behavior of the 'go' command-line tool (with some minor changes).
*/
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"
)

var Writer io.Writer = os.Stderr

type Program struct {
	Name     string    // program name
	Commands []Command // program commands
}

func (prog *Program) RunAndExit() {
	os.Exit(prog.Run(os.Args))
}

func (prog *Program) Run(args []string) (code int) {
	code = 1 // error code will only be 0 if the command successfully ran

	var templData = struct {
		Program *Program
		Command *Command
		Flags   []*flag.Flag
	}{
		Program: prog,
	}

	fs := flag.NewFlagSet(prog.Name, flag.ContinueOnError)
	fs.Usage = func() {}
	err := fs.Parse(args[1:])
	if (err == flag.ErrHelp) || (fs.NArg() == 0) || (fs.NArg() == 1 && fs.Arg(0) == "help") {
		// If '-help' or '-h', no command, or 'help' command print the program usage and exit
		Templ(programUsageTempl, &templData)
		return
	}

	for i := 0; i < len(prog.Commands); i++ {
		cmd := &prog.Commands[i]

		if fs.Arg(0) == cmd.Name || fs.Arg(1) == cmd.Name {
			templData.Command = cmd

			cfs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
			cfs.Usage = func() {
				Templ(commandUsageTempl, &templData)
			}
			cmd.SetupFlags(cfs)
			cfs.Parse(fs.Args()[1:])

			if fs.Arg(0) == "help" {
				templData.Flags = make([]*flag.Flag, 0)
				cfs.VisitAll(func(f *flag.Flag) {
					templData.Flags = append(templData.Flags, f)
				})
				Templ(commandHelpTempl, &templData)
				return

			} else if fs.Arg(0) == cmd.Name {
				cmd.Run(cfs.Args())
				code = 0 // if we reached here the command has successfully run
				return
			}
		}
	}

	Infof(`%s: unknown command "%s"
Run '%s help' for usage.
`, prog.Name, fs.Arg(1), prog.Name)

	return
}

type Command struct {
	Name string // name of the command

	ShortDescription string // short description of the command
	Description      string // long description of the command

	ShortUsage string

	SetupFlags func(fs *flag.FlagSet)
	Run        func(args []string)
}

func Info(s string) {
	fmt.Fprint(Writer, s)
}

func Infof(format string, args ...interface{}) {
	fmt.Fprintf(Writer, format, args...)
}

func Fatal(s string) {
	Info(s)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	Infof(format, args...)
	os.Exit(1)
}

func Templ(t *template.Template, data interface{}) {
	if err := t.Execute(Writer, data); err != nil {
		panic(err)
	}
}

var programUsageTempl = template.Must(template.New("").Parse(`Usage: {{$.Program.Name}} command [options]

Available commands:
{{range $.Program.Commands}}
    {{printf "%-10s %s" .Name .ShortDescription}}
{{- end}}

Use "{{$.Program.Name}} help [command]" for more Information about that command.
`))

var commandUsageTempl = template.Must(template.New("").Parse(`Usage: {{$.Program.Name}} {{$.Command.Name}} {{$.Command.ShortUsage}}
Run '{{$.Program.Name}} help {{$.Command.Name}}' for details.
`))

var commandHelpTempl = template.Must(template.New("").Parse(`Usage: {{$.Program.Name}} {{$.Command.Name}} {{$.Command.ShortUsage}}

{{$.Command.Description}}

Options:
{{range $.Flags}}
    -{{printf "%-10s %s (default=%s)" .Name .Usage .DefValue}}
{{- end}}
`))
