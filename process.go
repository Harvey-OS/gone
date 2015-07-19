package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"encoding/csv"
)

func main() {
	b, err := ioutil.ReadFile("PORT_TO_HARVEY.csv")
	if err != nil {
		log.Fatalf("%v", err)
	}
	r := bytes.NewReader(b)
	c := csv.NewReader(r)
	lines, err := c.ReadAll()
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", lines)
}
