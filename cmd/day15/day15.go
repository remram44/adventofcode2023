package main

import "bufio"
import "fmt"
import "log"
import "os"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day15.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Read the one line
	if !fileScanner.Scan() {
		log.Fatal("No line in file")
	}
	line := fileScanner.Text()
	if fileScanner.Scan() {
		log.Fatal("Multiple lines in file")
	}

	// Iterate on comma-separated steps
	start := 0
	sumOfHashes := 0
	for pos := 0; pos < len(line); pos += 1 {
		if line[pos] == ',' {
			sumOfHashes += int(hash(line[start:pos]))
			start = pos + 1
		}
	}
	if start < len(line) {
		sumOfHashes += int(hash(line[start:]))
	}

	fmt.Println(sumOfHashes)
}

func hash(s string) byte {
	var result byte = 0
	for i := 0; i < len(s); i += 1 {
		result += s[i]
		result *= 17
	}
	return result
}
