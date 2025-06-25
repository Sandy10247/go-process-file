package main

import (
	"log"
	"os"
	"unicode"

	"processfile/processor"
)

type Result struct {
	totalCapitals int
}

func (r *Result) CalculatTotalCapitals(data []byte) {
	totalCapitals := 0
	for _, v := range string(data) {
		if unicode.IsUpper(v) {
			totalCapitals++
		}
	}
	r.totalCapitals += totalCapitals
}

func main() {
	log.Println("hello World")

	file, err := os.Open("dummy_10_rows.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	res := &Result{}

	processor.ProcessFile(file, 10, 1_000, res.CalculatTotalCapitals)

	log.Println("totoal Capitals :- ", res.totalCapitals)
}
