package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"encoding/csv"
	"os"
)

var lines [][]string

func done(){
//	if err := ioutil.WriteFile("PORT_TO_HARVEY_NEXT.csv", []byte(lines[:][:]), 0666); err != nil {
//		log.Fatalf("%v", err)
//	}
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
		switch b {
			case "ignore":
				lines[i][1] = "ignore"
			case "skip":
			case "fixed":
				lines[i][1] = "fixed"
			case "edit":
				fmt.Print("Run an editor")
			case "copy":
				fmt.Printf("Copy %v somewhere", lines[i][0])
				lines[i][1] = "copy"
			case "exit":
				done()
			default:
				fmt.Printf("?")
		}


	}
	done()
}
