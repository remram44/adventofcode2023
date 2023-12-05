package main

import "bufio"
import "fmt"
import "log"
import "os"
import "sort"
import "strings"

import "github.com/remram44/adventofcode2023"

func main() {
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
	var numbers []int
	for pos < len(line) {
		var num int
		pos, num = aoc.ReadNumber(line, pos)
		numbers = append(numbers, num)
		pos += 1
	}

	// Sort seeds
	sort.Ints(numbers)

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
		newNumbers := make([]int, len(numbers))
		for i := 0; i < len(numbers); i += 1 {
			newNumbers[i] = -1
		}

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
				log.Fatalf("%v != %v", pos, len(entry))
			}

			// Change numbers via this map entry
			changeFrom := sort.SearchInts(numbers, src_start)
			changeTo := sort.SearchInts(numbers, src_start+length)
			for i := changeFrom; i < changeTo; i += 1 {
				// Only map once
				if newNumbers[i] != -1 {
					continue
				}

				newNumbers[i] = numbers[i] + dest_start - src_start
			}
		}

		// Copy the rest of the numbers unchanged
		for i := 0; i < len(numbers); i += 1 {
			if newNumbers[i] == -1 {
				newNumbers[i] = numbers[i]
			}
		}

		// Swap the arrays, sort it
		numbers = newNumbers
		sort.Ints(numbers)
	}

	fmt.Println(numbers[0])
}
