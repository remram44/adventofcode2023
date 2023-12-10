package main

import "bufio"
import "fmt"
import "log"
import "os"
import "slices"

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
	total1 := 0
	total2 := 0
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

		// Extrapolate forward
		total1 += extrapolateSeries(series, 1)

		// Extrapolate backwards
		slices.Reverse(series)
		total2 += extrapolateSeries(series, -1)
	}

	fmt.Println(total1)
	fmt.Println(total2)
}

func extrapolateSeries(series []int, sign int) int {
	fmt.Println(series)
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
			derivedSeries = append(derivedSeries, sign*(series[i+1]-series[i]))
		}

		extrapolated := extrapolateSeries(derivedSeries, sign)
		return series[len(series)-1] + sign*extrapolated
	}
}
