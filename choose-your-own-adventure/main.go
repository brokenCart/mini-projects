package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	filenamePtr := flag.String("file", "gopher.json", "a json file in a specific format for the story (default gopher.json)")
	flag.Parse()

	jsonFile, err := os.Open(*filenamePtr)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	stories, err := parseJSON(byteValue)

	if err != nil {
		return err
	}
	fmt.Println(stories["intro"])

	return nil
}
