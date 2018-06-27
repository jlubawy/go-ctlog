// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"flag"
	"fmt"
	"os"
	"text/template"
)

type Program struct {
	Name     string    // program name
	Commands []Command // program commands
}

func (prog *Program) Run() {
	var templData = struct {
		Program *Program
		Command *Command
		Flags   []*flag.Flag
	}{
		Program: prog,
	}

	usageFunc := func() {
		Templ(programUsageTempl, &templData)
		os.Exit(1)
	}

	fs := flag.NewFlagSet(prog.Name, flag.ExitOnError)
	fs.Usage = usageFunc
	fs.Parse(os.Args[1:])

	if (fs.NArg() == 0) || (fs.NArg() == 1 && fs.Arg(0) == "help") {
		// If no command or 'help' command print the program usage and exit
		usageFunc()
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
				os.Exit(1)

			} else if fs.Arg(0) == cmd.Name {
				cmd.Run(cfs.Args())
				os.Exit(0)
			}
		}
	}

	Fatalf(`%s: unknown command "%s"
Run '%s help' for usage.
`, prog.Name, fs.Arg(1), prog.Name)
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
	fmt.Fprint(os.Stderr, s)
}

func Infof(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
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
	if err := t.Execute(os.Stderr, data); err != nil {
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

var commandUsageTempl = template.Must(template.New("").Parse(`usage: {{$.Program.Name}} {{$.Command.Name}} {{$.Command.ShortUsage}}
Run '{{$.Program.Name}} help {{$.Command.Name}}' for details.
`))

var commandHelpTempl = template.Must(template.New("").Parse(`usage: {{$.Program.Name}} {{$.Command.Name}} {{$.Command.ShortUsage}}

{{$.Command.Description}}

Options:
{{range $.Flags}}
    -{{printf "%-10s %s (default=%s)" .Name .Usage .DefValue}}
{{- end}}
`))
