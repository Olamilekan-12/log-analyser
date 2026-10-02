package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type LogEntry struct {
	Timestamp string
	IP        string
	Level     string
	Message   string
}

type ErrorCount struct {
	Message string
	Count   int
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

func toSlice(counts map[string]int) []ErrorCount {
	var result []ErrorCount
	for message, count := range counts {
		result = append(result, ErrorCount{
			Message: message,
			Count:   count,
		})
	}
	return result
}

func countErrors(entries []LogEntry) map[string]int {
	counts := make(map[string]int)
	for _, entry := range entries {
		if entry.Level == "ERROR" {
			counts[entry.Message]++
		}
	}
	return counts
}

func readEntries(fp string) ([]LogEntry, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var entries []LogEntry
	for scanner.Scan() {
		entry, err := parseLine(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: loganalyzer <logfile>")
		os.Exit(1)
	}
	entries, err := readEntries(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading log: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(len(entries))
	errorCounts := countErrors(entries)
	errorSlice := toSlice(errorCounts)
	sort.Slice(errorSlice, func(i, j int) bool {
		if errorSlice[i].Count > errorSlice[j].Count {
			return errorSlice[i].Count > errorSlice[j].Count
		}
		return errorSlice[i].Message < errorSlice[j].Message
	})
	fmt.Println(errorSlice)
}
