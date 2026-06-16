package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"urlshortener"
)

func main() {
	filenamePtr := flag.String("file", "mappings.yml", "a json/yml file which has path and url as fields (default \"mappings.yml\")")

	db, err := urlshortener.GetDB("mappings.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = urlshortener.CreateTable(db)
	if err != nil {
		panic(err)
	}

	rows, err := urlshortener.GetData(db)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	mux := defaultMux()
	dbHandler, err := urlshortener.DBHandler(rows, mux)
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(*filenamePtr)
	if err != nil {
		panic(err)
	}

	extension := strings.Split(*filenamePtr, ".")[1]
	var fileHandler http.HandlerFunc

	switch extension {
	case "yml":
		fileHandler, err = urlshortener.YAMLHandler([]byte(content), dbHandler)
		if err != nil {
			panic(err)
		}
	case "json":
		fileHandler, err = urlshortener.JSONHandler([]byte(content), dbHandler)
		if err != nil {
			panic(err)
		}
	default:
		panic(errors.New("Invalid file extension: Only .yml and .json allowed"))
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", fileHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
