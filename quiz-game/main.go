package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file := openFile(*csvFileName)
	problems := getProblems(file)
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	answerCh := make(chan string)

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		go getAnswerFromUser(answerCh)

		select {
		case <-timer.C:
			printScore(correct, len(problems), true)
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	printScore(correct, len(problems), false)
}

type problem struct {
	q string
	a string
}

func openFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to open the CSV file: %s\n", fileName)
		exit(errMsg)
	}
	return file
}

func getProblems(file *os.File) []problem {
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}
	problems := parseLines(lines)
	return problems
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func getAnswerFromUser(answerCh chan string) {
	var answer string
	fmt.Scanf("%s\n", &answer)
	answerCh <- answer
}

func printScore(correct int, total int, newLine bool) {
	if newLine {
		fmt.Printf("\nYou scored %d out of %d\n", correct, total)
	} else {
		fmt.Printf("You scored %d out of %d\n", correct, total)
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
