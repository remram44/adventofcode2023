package main

import "bufio"
import "fmt"
import "log"
import "os"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day01.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	sumOfValues := 0

	// Iterate on lines
	for fileScanner.Scan() {
		sumOfValues += extractCalibrationValue(fileScanner.Text())
	}

	// Print sum
	fmt.Println(sumOfValues)
}

func extractCalibrationValue(line string) int {
	var firstDigit rune = 0
	var lastDigit rune = 0
	for _, char := range line {
		if '0' <= char && char <= '9' {
			if firstDigit == 0 {
				firstDigit = char
			}
			lastDigit = char
		}
	}
	return int(firstDigit-'0')*10 + int(lastDigit-'0')
}
