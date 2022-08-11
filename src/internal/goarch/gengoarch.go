// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	       "r2.is/gocompiler/vfs/os"
	"strconv"
	"strings"
)

var goarches []string

func main() {
	data, err := os.ReadFile("../../go/build/syslist.go")
	if err != nil {
		log.Fatal(err)
	}
	const goarchPrefix = `const goarchList = `
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, goarchPrefix) {
			text, err := strconv.Unquote(strings.TrimPrefix(line, goarchPrefix))
			if err != nil {
				log.Fatalf("parsing goarchList: %v", err)
			}
			goarches = strings.Fields(text)
		}
	}

	for _, target := range goarches {
		if target == "amd64p32" {
			continue
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "// Code generated by gengoarch.go using 'go generate'. DO NOT EDIT.\n\n")
		fmt.Fprintf(&buf, "//go:build %s\n", target) // must explicitly include target for bootstrapping purposes
		fmt.Fprintf(&buf, "package goarch\n\n")
		fmt.Fprintf(&buf, "const GOARCH = `%s`\n\n", target)
		for _, goarch := range goarches {
			value := 0
			if goarch == target {
				value = 1
			}
			fmt.Fprintf(&buf, "const Is%s = %d\n", strings.Title(goarch), value)
		}
		err := os.WriteFile("zgoarch_"+target+".go", buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
}
