package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"encoding/csv"
	"os"
	"strings"
)

var lines [][]string

func done(){
	c := csv.NewWriter(os.Stdout)
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
	fmt.Printf("%v\n", lines)
	for i := range lines {
		var b string
		if lines[i][1] != "" {
			continue
		}
		fmt.Printf("ASK ABOUT %v\n", lines[i])
		n, err := fmt.Scanln(&b)
		if n == 0 {
			break
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		repeat:
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
					editor = "/bin/vi"
				}
				fmt.Printf("%v %v\n", editor, lines[i][0])
				break repeat
			case "copy":
				fmt.Printf("Copy %v ", lines[i][0])
				b, err := ioutil.ReadFile(lines[i][0])
				if err != nil {
					log.Printf("%v", err)
					done()
				}
				if ! strings.Contains(lines[i][0], "plan9.go") {
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


	}
	done()
}
