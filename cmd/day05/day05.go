package main

import "bufio"
import "fmt"
import "log"
import "os"
import "slices"
import "sort"
import "strings"

import "github.com/remram44/adventofcode2023"

type NumberRange struct {
	start  int
	length int
}

func SortNumberRangeArray(arr []NumberRange) {
	sort.Slice(arr, func(i int, j int) bool {
		return arr[i].start < arr[j].start
	})
}

func main() {
	parseAlmanac(false)
	parseAlmanac(true)
}

func parseAlmanac(seedsAreRanges bool) {
	// Open the input file
	readFile, err := os.Open("inputs/day05.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Read seeds into the "numbers" array
	if !fileScanner.Scan() {
	}
	line := fileScanner.Text()
	if !strings.HasPrefix(line, "seeds: ") {
		log.Fatalf("Invalid seeds line: %v", line)
	}
	pos := 7
	var numbers []NumberRange
	if !seedsAreRanges {
		// Part 1: each number is a single seed number
		for pos < len(line) {
			var num int
			pos, num = aoc.ReadNumber(line, pos)
			numbers = append(numbers, NumberRange{
				start:  num,
				length: 1,
			})
			pos += 1
		}
	} else {
		// Part 2: each pair of numbers is a range
		for pos < len(line) {
			var start int
			var length int
			pos, start = aoc.ReadNumber(line, pos)
			pos += 1
			pos, length = aoc.ReadNumber(line, pos)
			numbers = append(numbers, NumberRange{
				start:  start,
				length: length,
			})
			pos += 1
		}
	}

	// Sort seeds
	SortNumberRangeArray(numbers)

	if !fileScanner.Scan() || fileScanner.Text() != "" {
		log.Fatalf("Missing separator")
	}

	// Iterate on maps
	for fileScanner.Scan() {
		header := fileScanner.Text()
		if !strings.HasSuffix(header, " map:") {
			log.Fatalf("Invalid map header: %v", header)
		}

		// Allocate an array for the mapped numbers
		var newNumbers []NumberRange

		// Read map entries
		for fileScanner.Scan() {
			entry := fileScanner.Text()
			if entry == "" {
				break
			}

			pos := 0
			var dest_start, src_start, length int
			pos, dest_start = aoc.ReadNumber(entry, pos)
			pos += 1
			pos, src_start = aoc.ReadNumber(entry, pos)
			pos += 1
			pos, length = aoc.ReadNumber(entry, pos)
			if pos != len(entry) {
				log.Fatalf("Too much text on map entry line: %v", entry)
			}

			// Change numbers via this map entry
			changeFrom := sort.Search(len(numbers), func(i int) bool {
				return src_start < numbers[i].start+numbers[i].length
			})
			changeTo := sort.Search(len(numbers), func(i int) bool {
				return numbers[i].start >= src_start+length
			})
			for i := changeFrom; i < changeTo; i += 1 {
				// Only map once
				if numbers[i].length == 0 {
					continue
				}

				var leftOver []NumberRange

				// Add the first part of the range, that isn't getting mapped
				if numbers[i].start < src_start {
					leftOver = append(leftOver, NumberRange{
						start:  numbers[i].start,
						length: src_start - numbers[i].start,
					})
					numbers[i].length -= src_start - numbers[i].start
					numbers[i].start = src_start
				}

				// Add the last part of the range, that isn't getting mapped
				if numbers[i].start+numbers[i].length > src_start+length {
					leftOver = append(leftOver, NumberRange{
						start:  src_start + length,
						length: numbers[i].start + numbers[i].length - src_start - length,
					})
					numbers[i].length = src_start + length - numbers[i].start
				}

				// Map the range over
				newNumbers = append(newNumbers, NumberRange{
					start:  numbers[i].start + dest_start - src_start,
					length: numbers[i].length,
				})

				// Erase that range
				numbers[i].length = 0

				// Insert the left-over ranges in sorted order
				for _, number := range leftOver {
					idx := sort.Search(len(numbers), func(i int) bool {
						return number.start <= numbers[i].start
					})
					numbers = slices.Insert(numbers, idx, number)
					if idx <= i {
						i += 1
					}
					if idx < changeTo {
						changeTo += 1
					}
				}
			}
		}

		// Copy the rest of the ranges unchanged
		for i := 0; i < len(numbers); i += 1 {
			if numbers[i].length != 0 {
				newNumbers = append(newNumbers, numbers[i])
			}
		}

		// Swap the arrays, sort it
		numbers = newNumbers
		SortNumberRangeArray(numbers)
	}

	fmt.Println(numbers[0].start)
}
