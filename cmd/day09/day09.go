package main

import "bufio"
import "fmt"
import "log"
import "os"

import "github.com/remram44/adventofcode2023"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day09.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines
	total := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()

		// Read numbers
		pos := 0
		var series []int
		for pos < len(line) {
			pos = aoc.ReadSpaces(line, pos)
			var number int
			pos, number = aoc.ReadNumber(line, pos)
			series = append(series, number)
		}

		// Extrapolate
		total += extrapolateSeries(series)
	}

	fmt.Println(total)
}

func extrapolateSeries(series []int) int {
	allZeroes := true
	for _, elem := range series {
		if elem != 0 {
			allZeroes = false
			break
		}
	}

	if allZeroes {
		return 0
	} else {
		var derivedSeries []int
		for i := 0; i < len(series)-1; i += 1 {
			derivedSeries = append(derivedSeries, series[i+1]-series[i])
		}

		extrapolated := extrapolateSeries(derivedSeries)
		return series[len(series)-1] + extrapolated
	}
}
