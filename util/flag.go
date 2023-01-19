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

package util

import (
	"flag"
	"os"
	"strings"
	"text/template"
)

// PrintGlobalUsage prints a help message to standard error using the default
// flag set. See PrintUsage for details.
func PrintGlobalUsage(usage string) {
	PrintUsage(flag.CommandLine, usage)
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
func PrintUsage(flags *flag.FlagSet, usage string) {
	usage = strings.TrimSpace(usage) + "\n"
	t := template.Must(template.New("").Parse(usage))
	t.Execute(os.Stderr, map[string]interface{}{
		"Program": Program(),
		"Name":    flags.Name(),
		"PrintDefaults": func() string {
			var b strings.Builder
			flags.SetOutput(&b)
			flags.PrintDefaults()
			return strings.TrimSpace(b.String())
		},
	})
}
