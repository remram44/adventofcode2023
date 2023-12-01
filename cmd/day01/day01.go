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
	var firstDigit int = -1
	var lastDigit int = -1
	for _, char := range line {
		if '0' <= char && char <= '9' {
			digit := int(char - '0')
			if firstDigit == -1 {
				firstDigit = digit
			}
			lastDigit = digit
		}
	}
	return firstDigit*10 + lastDigit
}
