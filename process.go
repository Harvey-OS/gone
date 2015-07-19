package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var lines [][]string

var debug = flag.Bool("debug", false, "Debug output")

func done() {
	o, err := os.Create("PORT_TO_HARVEY_NEXT.csv")
	if err != nil {
		log.Fatalf("%v", err)
	}
	c := csv.NewWriter(o)
	if err := c.WriteAll(lines); err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(0)
}

func main() {
	b, err := ioutil.ReadFile("PORT_TO_HARVEY.csv")
	if err != nil {
		log.Fatalf("%v", err)
	}
	r := bytes.NewReader(b)
	c := csv.NewReader(r)
	lines, err = c.ReadAll()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if *debug {
		fmt.Printf("%v\n", lines)
	}
	for i := range lines {
		var b string
		if lines[i][1] != "" {
			continue
		}
		for {
			advance := true
			fmt.Printf("ASK ABOUT %v\n", lines[i])
			n, err := fmt.Scanln(&b)
			if n == 0 {
				break
			}
			if err != nil {
				log.Fatalf("%v", err)
			}
			switch b {
			case "ignore":
				lines[i][1] = "ignore"
			case "skip":
			case "fixed":
				lines[i][1] = "fixed"
			case "edit":
				fmt.Print("Run an editor")
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "/usr/bin/vi"
				}
				fmt.Printf("%v %v\n", editor, lines[i][0])
				cmd := exec.Command(editor, lines[i][0])
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					log.Printf("%v", err)
					done()
				}
				advance = false
			case "copy":
				fmt.Printf("Copy %v ", lines[i][0])
				b, err := ioutil.ReadFile(lines[i][0])
				if err != nil {
					log.Printf("%v", err)
					done()
				}
				if !strings.Contains(lines[i][0], "plan9.go") {
					log.Printf("Can't split %v around plan9 9", lines[i][0])
					done()
				}
				n := strings.SplitN(lines[i][0], "plan9.go", 2)
				newname := n[0] + "harvey.go"
				fmt.Printf("to %v ", newname)
				if err := ioutil.WriteFile(newname, b, 0666); err != nil {
					log.Printf("%v", err)
					done()
				}
				lines[i][1] = "copy"
			case "exit":
				done()
			default:
				fmt.Printf("?")
			}
			if advance {
				break
			}
		}

	}
	done()
}
