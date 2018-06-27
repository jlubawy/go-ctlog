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

// A Program is named command-line program with a collection of supported
// sub-commands.
type Program struct {
	Name     string    // name of the programs
	Commands []Command // supported commands
}

// A Command is a sub-command supported by a given program.
type Command struct {
	Name             string // name of the command
	ShortDescription string // short description of the command
	Description      string // long description of the command
	ShortUsage       string // short usage description

	// SetupFlags is called before flag.Parse allowing the command to setup
	// any options it needs.
	SetupFlags func(fs *flag.FlagSet)

	// Run is the entry point to the command, it is passed the arguments after
	// the command name.
	Run func(args []string)
}

// RunAndExit is the same as Run but it calls os.Exit with the error code
// returned by Run.
func (prog *Program) RunAndExit() {
	os.Exit(prog.Run(os.Args))
}

// Run runs the program or prints help messages if requested, returning 0 if
// successful, DefaultErrorCode if unsuccessful, or it may not return at all
// if the command exited the program. Typically RunAndExit should be used
// instead.
func (prog *Program) Run(args []string) (code int) {
	code = DefaultErrorCode

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

// Writer is the io.Writer that the Info, Infof, Fatal, and Fatalf functions
// write to. It defaults to os.Stderr.
var Writer io.Writer = os.Stderr

// DefaultErrorCode is the default exit code used by all calls to os.Exit in
// this package.
var DefaultErrorCode = 1

// Info logs a message to Writer using fmt.Fprint.
func Info(s string) {
	fmt.Fprint(Writer, s)
}

// Infof logs a message to Writer using fmt.Fprintf.
func Infof(format string, args ...interface{}) {
	fmt.Fprintf(Writer, format, args...)
}

// Fatal is similar to Info but it calls os.Exit(DefaultErrorCode) after.
func Fatal(s string) {
	Info(s)
	os.Exit(DefaultErrorCode)
}

// Fatalf is similar to Infof but it calls os.Exit(DefaultErrorCode) after.
func Fatalf(format string, args ...interface{}) {
	Infof(format, args...)
	os.Exit(DefaultErrorCode)
}

// Templ executes the provided template writing to Writer, calling panic if
// there was any error.
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
