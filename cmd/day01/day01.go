package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strings"

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

	sumOfValuesPart1 := 0
	sumOfValuesPart2 := 0

	// Iterate on lines
	for fileScanner.Scan() {
		sumOfValuesPart1 += extractCalibrationValueFromDigits(fileScanner.Text())
		sumOfValuesPart2 += extractCalibrationValueFromDigitsAndLetters(fileScanner.Text())
	}

	// Print sum
	fmt.Println(sumOfValuesPart1)
	fmt.Println(sumOfValuesPart2)
}

func extractCalibrationValueFromDigits(line string) int {
	var firstDigit = -1
	var lastDigit = -1
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

func extractCalibrationValueFromDigitsAndLetters(line string) int {
	var firstDigit = -1
	var lastDigit = -1
	for pos, char := range line {
		digit := -1
		if '0' <= char && char <= '9' {
			digit = int(char - '0')
		} else if strings.HasPrefix(line[pos:], "one") {
			digit = 1
		} else if strings.HasPrefix(line[pos:], "two") {
			digit = 2
		} else if strings.HasPrefix(line[pos:], "three") {
			digit = 3
		} else if strings.HasPrefix(line[pos:], "four") {
			digit = 4
		} else if strings.HasPrefix(line[pos:], "five") {
			digit = 5
		} else if strings.HasPrefix(line[pos:], "six") {
			digit = 6
		} else if strings.HasPrefix(line[pos:], "seven") {
			digit = 7
		} else if strings.HasPrefix(line[pos:], "eight") {
			digit = 8
		} else if strings.HasPrefix(line[pos:], "nine") {
			digit = 9
		}

		if digit != -1 {
			if firstDigit == -1 {
				firstDigit = digit
			}
			lastDigit = digit
		}
	}
	return firstDigit*10 + lastDigit
}
