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

// Package generate provides functions for implementing generators.
package generate

import (
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sstallion/go-tools/util"
)

// Args returns the command line arguments, starting with the base program
// name as a string.
func Args() string {
	return strings.Join(util.Args(), " ")
}

// GoArch returns the value of the $GOARCH environment variable, or panics if
// the variable is not present.
func GoArch() string {
	return util.MustEnv("GOARCH")
}

// GoOs returns the value of the $GOOS environment variable, or panics if the
// variable is not present.
func GoOs() string {
	return util.MustEnv("GOOS")
}

// GoFile returns the value of the $GOFILE environment variable, or panics if
// the variable is not present.
func GoFile() string {
	return util.MustEnv("GOFILE")
}

// GoLine returns the value of the $GOLINE environment variable, or panics if
// the variable is not present.
func GoLine() string {
	return util.MustEnv("GOLINE")
}

// GoPackage returns the value of the $GOPACKAGE environment variable, or
// panics if the variable is not present.
func GoPackage() string {
	return util.MustEnv("GOPACKAGE")
}

// GoRoot returns the value of the $GOROOT environment variable, or panics if
// the variable is not present.
func GoRoot() string {
	return util.MustEnv("GOROOT")
}

// Dollar returns the value of the $DOLLAR environment variable, or panics if
// the variable is not present.
func Dollar() string {
	return util.MustEnv("DOLLAR")
}

// GoCmd returns an *exec.Cmd to execute "$GOROOT/bin/go" with the given
// arguments.
func GoCmd(arguments []string) *exec.Cmd {
	name := filepath.Join(GoRoot(), "bin", "go")
	return exec.Command(name, arguments...)
}

// GoCmd returns an *exec.Cmd to execute "$GOROOT/bin/go run" with the given
// arguments.
func GoRunCmd(pkg string, arguments []string) *exec.Cmd {
	return GoCmd(append([]string{"run", pkg}, arguments...))
}

// GoCmd returns an *exec.Cmd to execute "$GOROOT/bin/go generate" with the
// given arguments.
func GoGenerateCmd(arguments []string) *exec.Cmd {
	return GoCmd(append([]string{"generate"}, arguments...))
}

// WriteFile wraps os.WriteFile to use standard output if name is "-".
func WriteFile(name string, data []byte, perm os.FileMode) error {
	if name == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(name, data, perm)
}

// WriteSource first passes data through format.Source before calling
// os.Writefile. Standard output will be used if name is "-".
func WriteSource(name string, data []byte, perm os.FileMode) error {
	src, err := format.Source(data)
	if err != nil {
		return err
	}
	return WriteFile(name, src, perm)
}
