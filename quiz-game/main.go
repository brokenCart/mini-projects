package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Declare and Parse flags
	problemsCSVFile := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer' (default \"problems.csv\")")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds (default 30)")
	randomOrder := flag.Bool("random", false, "the questions appear in a random order if this flag is used")
	flag.Parse()

	// Open and Read the problems file
	file, err := os.Open(*problemsCSVFile)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	// To close the file when the function returns
	defer file.Close()

	// Parse the CSV file into problems
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = 2 // Set the number of columns (ReadAll() raises error if 2 fields are not parsed per row)

	lines, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("Failed to read CSV: %w", err)
	}
	problems := parseLines(lines)

	// Shuffle the problems if the -random flag is used
	if *randomOrder {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	// Format Width for the question number
	formatWidth := int(math.Log10(float64(len(problems)))) + 1

	// Initializing the input reader
	inputReader := bufio.NewReader(os.Stdin)

	// Starting the timer when the user presses the enter key
	fmt.Printf("Number of Questions = %d\nTime = %ds\nPress Enter key to start the quiz...", len(problems), *timeLimit)
	inputReader.ReadString('\n')
	timer := time.NewTimer(time.Second * time.Duration(*timeLimit))

	// Go through each question and calculate the score
	correctQuestions := 0
	for i, prob := range problems {
		// Prompt the user for an answer
		fmt.Printf("%*d. %s = ", formatWidth, i+1, prob.question)

		// Create a channel to get the input
		answerCh := make(chan inputResult, 1)

		// Create a separate goroutine for getting input
		// so that it doesn't blocks the select statement
		go func() {
			input, err := inputReader.ReadString('\n')
			inputRes := inputResult{
				input,
				err,
			}
			answerCh <- inputRes
		}()

		// Check whichever channel got the data first
		// and process accordingly
		select {
		// The timer channel got the data first,
		// this means that the time has ran out.
		case <-timer.C:
			fmt.Printf("\n\nTime's up!\nYou scored %d out of %d.\n", correctQuestions, len(problems))
			return nil
		// The answer channel got the data first,
		// this means that the input was sent by the user
		case answerRes := <-answerCh:
			// Check if there was an error in reading the input
			if answerRes.err != nil {
				return fmt.Errorf("Error reading input: %w", answerRes.err)
			}

			// Check if the answer was correct
			input := strings.ToLower(strings.TrimSpace(answerRes.input))
			if input == prob.answer {
				correctQuestions++
			}
		}
	}

	fmt.Printf("\nYou completed the quiz in time!\nYou scored %d out of %d.\n", correctQuestions, len(problems))
	return nil
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, row := range lines {
		// Get the question and answer
		question, answer := row[0], row[1]
		answer = strings.ToLower(strings.TrimSpace(answer))
		problems[i] = problem{
			question,
			answer,
		}
	}
	return problems
}

type problem struct {
	question string
	answer   string
}

type inputResult struct {
	input string
	err   error
}
