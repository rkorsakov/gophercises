package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type quiz struct {
	question string
	answer   string
}

type config struct {
	fileName string
	duration time.Duration
}

func setUpFlags() config {
	filePtr := flag.String("f", "problems.csv", "CSV file to read from")
	durationPtr := flag.Int("t", 3, "Duration of time to answer all questions")
	flag.Parse()
	duration := time.Duration(*durationPtr) * time.Second
	return config{fileName: *filePtr, duration: duration}
}

func main() {
	cfg := setUpFlags()
	file, err := os.Open(cfg.fileName)
	defer file.Close()
	if err != nil {
		shutdown(err)
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		shutdown(err)
	}
	parsed, err := parse(records)
	length := len(parsed)
	var rightAnswers int
	fmt.Printf("Press enter to start a game! You have %s seconds to answer", cfg.duration.String())
	fmt.Scanln()
	timer := time.NewTimer(cfg.duration)
	defer timer.Stop()
	for i, _ := range parsed {
		fmt.Println("Question: ", parsed[i].question)
		input := make(chan string)
		go func() {
			var guess string
			fmt.Scanln(&guess)
			input <- guess
		}()
		fmt.Println("Say your answer: ")
		select {
		case <-timer.C:
			fmt.Println("\nTime's up")
			fmt.Printf("Right answers: %d from %d", rightAnswers, length)
			return
		case answer := <-input:
			if answer == parsed[i].answer {
				rightAnswers++
			}
		}
	}
	fmt.Printf("Right answers: %d from %d", rightAnswers, length)
}

func parse(records [][]string) ([]quiz, error) {
	var returnSlice []quiz
	for _, record := range records {
		returnSlice = append(returnSlice, quiz{question: record[0], answer: record[1]})
	}
	return returnSlice, nil
}

func shutdown(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}
