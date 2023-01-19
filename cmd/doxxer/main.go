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

//go:generate doxxer . -h
package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/sstallion/go-tools/generate"
	"github.com/sstallion/go-tools/util"
)

const tmplText = `
// Code generated by "{{ .Args }}"; DO NOT EDIT.

/*
{{ range .Text -}}
{{ . }}
{{ end -}}
*/
package {{ .Package }}
`

var output string

func usage() {
	util.PrintGlobalUsage(`
Doxxer is a tool that generates documentation for command line applications.

Usage:

  {{ .Program }} [-o output] <package> [arguments...]

Flags:

  {{ call .PrintDefaults }}

Typically, arguments are specified using "//go:generate" directives, which are
passed verbatim to "go run" to generate output. Output is then passed through
"gofmt" and finally written to a file, which by default is doc.go.

The following example demonstrates generating documentation for an application
that makes use of the standard flag package:

Example:

  //go:generate doxxer . -h
  package main

If doxxer is called directly from the command line, the $GOROOT and $GOPACKAGE
environment variables must be defined as documented in "go help generate".

Report issues to https://github.com/sstallion/go-tools/issues.
`)
}

func main() {
	// Intercept panics due to missing environment variables in the event
	// we're called from the command line directly:
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	log.SetFlags(0)
	log.SetPrefix("doxxer: ")

	flag.Usage = usage
	flag.StringVar(&output, "o", "doc.go", "Write `output` to file")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	var in bytes.Buffer
	cmd := generate.GoRunCmd(args[0], args[1:])
	cmd.Stdout = &in
	cmd.Stderr = &in
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	text := make(chan string)
	go func() {
		scanner := bufio.NewScanner(&in)
		for scanner.Scan() {
			text <- scanner.Text()
		}
		close(text)
	}()

	var out bytes.Buffer
	t := template.Must(template.New("").Parse(tmplText))
	t.Execute(&out, map[string]interface{}{
		"Args":    generate.Args(),
		"Text":    text,
		"Package": generate.GoPackage(),
	})
	if err := generate.WriteSource(output, out.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
}