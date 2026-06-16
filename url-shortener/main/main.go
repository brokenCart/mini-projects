package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"urlshortener"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Initialized the flags
	filenamePtr := flag.String("file", "mappings.yml", "a json/yml file which has path and url as fields (default \"mappings.yml\")")
	portPtr := flag.String("port", "8080", "the port to start the web server on (default 8080)")
	flag.Parse()

	// Get the db connection
	db, err := urlshortener.GetDB("mappings.db")
	if err != nil {
		return err
	}
	defer db.Close() // Close the db connection when the function exits

	// Create the table if it doesn't exist
	err = urlshortener.CreateTable(db)
	if err != nil {
		return err
	}

	// Get the data from the db
	rows, err := urlshortener.GetData(db)
	if err != nil {
		return err
	}
	defer rows.Close() // Close the rows when the function exits

	// Create the default mux
	mux := defaultMux()

	// Create the db handler
	dbHandler, err := urlshortener.DBHandler(rows, mux)
	if err != nil {
		return err
	}

	// Read the json/yaml file
	content, err := os.ReadFile(*filenamePtr)
	if err != nil {
		return err
	}

	// Get the file extension
	extension := filepath.Ext(*filenamePtr)
	var fileHandler http.HandlerFunc

	// Create the file handler based on the file extension
	switch extension {
	case "yml":
		fileHandler, err = urlshortener.YAMLHandler([]byte(content), dbHandler)
		if err != nil {
			return err
		}
	case "json":
		fileHandler, err = urlshortener.JSONHandler([]byte(content), dbHandler)
		if err != nil {
			return err
		}
	default:
		return errors.New("Invalid file extension: Only .yml and .json allowed")
	}

	// Start the server
	fmt.Println("Starting the server on :" + *portPtr)
	if err := http.ListenAndServe(":"+*portPtr, fileHandler); err != nil {
		return err
	}
	return nil
}

// defaultMux returns a default ServeMux with a simple hello handler.
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

// hello is a simple handler that writes "Hello, world!" to the response.
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
