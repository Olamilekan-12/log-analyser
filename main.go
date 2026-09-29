package main

import (
	"bufio"
	"fmt"
	"os"
)

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
