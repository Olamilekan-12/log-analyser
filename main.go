package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type LogEntry struct {
	Timestamp string
	IP        string
	Level     string
	Message   string
}

func parseLine(line string) (LogEntry, error) {
	linesSlice := strings.Fields(line)
	if len(linesSlice) < 4 {
		return LogEntry{}, fmt.Errorf("malformed line: %q", line)
	}
	messageExtraction := strings.Join(linesSlice[3:], " ")
	return LogEntry{
		Timestamp: linesSlice[0],
		IP:        linesSlice[1],
		Level:     linesSlice[2],
		Message:   messageExtraction,
	}, nil
}

func countLine(fp string) (int, error) {
	file, err := os.Open(fp)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lineCount int
	for scanner.Scan() {
		lineCount++
		entry, err := parseLine(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		fmt.Println(entry)
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return lineCount, nil
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: loganalyzer <logfile>")
		os.Exit(1)
	}
	lines, err := countLine(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting file %v\n", err)
		os.Exit(1)
	}
	fmt.Println(lines)
}
