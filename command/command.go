// Copyright (c) 2023 Steven Stallion <sstallion@gmail.com>
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

// Package command adds support for command processing to the standard flag
// package.
//
// # Usage
//
// Commands can be thought of as a wrapper for flag.FlagSet with some
// additional semantics. To get started, create a type that satisfies the
// Command interface and embeds a *flag.FlagSet to which additional flags may
// be bound:
//
//	type exampleCmd struct{
//		flags *flag.FlagSet
//		// additional flags
//	}
//
// Once defined, initialize the command and add it to the default set:
//
//	func init() {
//		cmd := &exampleCmd{flags: flag.NewFlagSet("example", flag.ExitOnError)}
//		cmd.flags.Usage = cmd.Usage
//		// bind additional flags
//		command.Add(cmd)
//	}
//
// Finally, after all commands have been defined, call:
//
//	command.Parse()
package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/sstallion/go-tools/util"
)

// ErrNArg indicates an incorrect number of arguments was passed to a command.
var ErrNArg = errors.New("wrong number of arguments")

// Command is an interface that allows an arbitrary type to be expressed as a
// command argument.
//
// Name returns the name of the command, which is used locate the command when
// parsing command line arguments. Description returns a short description of
// the command, which is used when displaying usage. If Description returns
// the empty string, the command will be considered unlisted and will not
// appear when printing commands. Parse parses the argument list and returns
// an error, if any. Run executes the command and returns an error, if any.
type Command interface {
	Name() string
	Description() string
	Usage()
	Parse(arguments []string) error
	Run() error
}

// CommandSet describes a set of defined commands.
type CommandSet []Command

// Add appends cmd to the set. If the command already exists in the set, it
// will be silently ignored.
func (cmds *CommandSet) Add(cmd Command) {
	if cmds.Lookup(cmd.Name()) != nil {
		return
	}
	*cmds = append(*cmds, cmd)
}

// Lookup returns the command from the set with the specified name. If no such
// command exists, nil is returned.
func (cmds *CommandSet) Lookup(name string) Command {
	for _, cmd := range *cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

// Visit visits all commands in insertion order, calling fn for each.
func (cmds *CommandSet) Visit(fn func(Command)) {
	for _, cmd := range *cmds {
		fn(cmd)
	}
}

// Parse parses the argument list and calls os.Exit with an appropriate error
// code, if any.
func (cmds *CommandSet) Parse(flags *flag.FlagSet, arguments []string) {
	flags.Parse(arguments)

	args := flags.Args()
	if len(args) > 0 {
		for _, cmd := range *cmds {
			if cmd.Name() != args[0] {
				continue
			}
			if err := cmd.Parse(args[1:]); err != nil {
				fmt.Fprintln(os.Stderr, err)
				cmd.Usage()
				os.Exit(1)
			}
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		fmt.Fprintf(os.Stderr, "invalid command: %s\n", args[0])
	}
	flags.Usage()
	os.Exit(1)
}

// PrintUsage prints a help message to standard error. Usage is typically a
// template string that can reference the following variables and functions:
//
// Variables:
//
//	Program
//		The base program name.
//	Name
//		The name of the flag set.
//
// Functions:
//
//	PrintDefaults
//		PrintDefaults prints to standard error the default values of
//		all defined command line flags in the flag set.
//	PrintCommands
//		PrintCommands prints to standard error the names and
//		descriptions of all defined commands in the command set.
func (cmds *CommandSet) PrintUsage(flags *flag.FlagSet, usage string) {
	usage = strings.TrimSpace(usage) + "\n"
	t := template.Must(template.New("").Parse(usage))
	t.Execute(os.Stderr, map[string]interface{}{
		"Program": util.Program(),
		"Name":    flags.Name(),
		"PrintDefaults": func() string {
			var b strings.Builder
			flags.SetOutput(&b)
			flags.PrintDefaults()
			return strings.TrimSpace(b.String())
		},
		"PrintCommands": func() string {
			var b strings.Builder
			w := tabwriter.NewWriter(&b, 2*8, 0, 0, ' ', 0)
			cmds.Visit(func(cmd Command) {
				if desc := cmd.Description(); desc != "" {
					fmt.Fprintf(w, "  %s\t%s\f", cmd.Name(), desc)
				}
			})
			return strings.TrimSpace(b.String())
		},
	})
}

// CommandLine is the default set of commands, parsed from os.Arg.
var CommandLine CommandSet

// Add appends cmd to the default set. If the command already exists in the
// set, it will be silently ignored.
func Add(cmd Command) {
	CommandLine.Add(cmd)
}

// Lookup returns the command from the default set with the specified name. If
// no such command exists, nil is returned.
func Lookup(name string) Command {
	return CommandLine.Lookup(name)
}

// Visit visits all commands in insertion order, calling fn for each.
func Visit(fn func(Command)) {
	CommandLine.Visit(fn)
}

// Parse parses the command line flags and commands from os.Args[1:].
func Parse() {
	CommandLine.Parse(flag.CommandLine, os.Args[1:])
}

// PrintGlobalUsage prints a help message to standard error using the default
// flag and command sets. See CommandSet.PrintUsage for details.
func PrintGlobalUsage(usage string) {
	PrintUsage(flag.CommandLine, usage)
}

// PrintUsage prints a help message to standard error using the default
// command set. See CommandSet.PrintUsage for details.
func PrintUsage(flags *flag.FlagSet, usage string) {
	CommandLine.PrintUsage(flags, usage)
}
