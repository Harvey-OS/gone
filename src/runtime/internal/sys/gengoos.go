// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var gooses, goarches []string

func main() {
	data, err := ioutil.ReadFile("../../../go/build/syslist.go")
	if err != nil {
		log.Fatal(err)
	}
	const (
		goosPrefix   = `const goosList = `
		goarchPrefix = `const goarchList = `
	)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, goosPrefix) {
			text, err := strconv.Unquote(strings.TrimPrefix(line, goosPrefix))
			if err != nil {
				log.Fatalf("parsing goosList %#q: %v", strings.TrimPrefix(line, goosPrefix), err)
			}
			gooses = strings.Fields(text)
		}
		if strings.HasPrefix(line, goarchPrefix) {
			text, err := strconv.Unquote(strings.TrimPrefix(line, goarchPrefix))
			if err != nil {
				log.Fatal("parsing goarchList: %v", err)
			}
			goarches = strings.Fields(text)
		}
	}

	for _, target := range gooses {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "// generated by gengoos.go using 'go generate'\n\n")
		if target == "linux" {
			fmt.Fprintf(&buf, "// +build !android\n\n") // must explicitly exclude android for linux
		}
		if target == "plan9" {
			fmt.Fprintf(&buf, "// +build !harvey\n\n") // must explicitly exclude android for linux
		}
		fmt.Fprintf(&buf, "package sys\n\n")
		fmt.Fprintf(&buf, "const GOOS = `%s`\n\n", target)
		for _, goos := range gooses {
			value := 0
			if goos == target {
				value = 1
			}
			fmt.Fprintf(&buf, "const Goos%s = %d\n", strings.Title(goos), value)
		}
		err := ioutil.WriteFile("zgoos_"+target+".go", buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, target := range goarches {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "// generated by gengoos.go using 'go generate'\n\n")
		fmt.Fprintf(&buf, "package sys\n\n")
		fmt.Fprintf(&buf, "const GOARCH = `%s`\n\n", target)
		for _, goarch := range goarches {
			value := 0
			if goarch == target {
				value = 1
			}
			fmt.Fprintf(&buf, "const Goarch%s = %d\n", strings.Title(goarch), value)
		}
		err := ioutil.WriteFile("zgoarch_"+target+".go", buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
}
