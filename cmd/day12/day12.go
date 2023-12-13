package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strings"

import "github.com/remram44/adventofcode2023"

type spring int

const (
	Empty spring = iota
	Present
	Unknown
)

func (s spring) String() string {
	switch s {
	case Empty:
		return "."
	case Present:
		return "#"
	case Unknown:
		return "?"
	default:
		return "!"
	}
}

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day12.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	sumOfArrangements := 0

	// Iterate on lines
	for fileScanner.Scan() {
		line := fileScanner.Text()

		// Read sequence of empty/present/unknown springs
		var sequence []spring
		pos := 0
		for pos < len(line) && line[pos] != ' ' {
			switch line[pos] {
			case '.':
				sequence = append(sequence, Empty)
			case '#':
				sequence = append(sequence, Present)
			case '?':
				sequence = append(sequence, Unknown)
			case ' ':
			default:
				log.Fatalf("Unexpected character %v", line[pos])
			}
			pos += 1
		}

		// Read list of groups
		var groupSizes []int
		for pos < len(line) {
			if line[pos] != ' ' && line[pos] != ',' {
				log.Fatalf("Invalid separator %v", line[pos])
			}
			pos += 1
			var num int
			pos, num = aoc.ReadNumber(line, pos)
			groupSizes = append(groupSizes, num)
		}

		// Find arrangements
		log.Print(line)
		arrangements := countArrangements(sequence, groupSizes, 0)
		log.Printf("%v -> %v", line, arrangements)
		log.Print()

		sumOfArrangements += arrangements
	}

	fmt.Println(sumOfArrangements)
}

func countArrangements(sequence []spring, groupSizes []int, depth int) int {
	log.Printf("%vcountArrangements(%v, %v)", strings.Repeat(" ", depth), sequence, groupSizes)
	if len(sequence) == 0 {
		if len(groupSizes) == 0 {
			log.Printf("%v done, return 1", strings.Repeat(" ", depth))
			return 1
		} else {
			log.Printf("%v fail, return 0", strings.Repeat(" ", depth))
			return 0
		}
	}

	switch sequence[0] {
	case Empty:
		log.Printf("%v empty, consume sequence", strings.Repeat(" ", depth))
		result := countArrangements(sequence[1:], groupSizes, depth+1)
		log.Printf("%v return %v", strings.Repeat(" ", depth), result)
		return result
	case Present:
		result := countArrangementsIfPresent(sequence, groupSizes, depth)
		log.Printf("%v return %v", strings.Repeat(" ", depth), result)
		return result
	case Unknown:
		total := 0
		log.Printf("%v try guess absent", strings.Repeat(" ", depth))
		total += countArrangements(sequence[1:], groupSizes, depth+1)
		log.Printf("%v try guess present", strings.Repeat(" ", depth))
		total += countArrangementsIfPresent(sequence, groupSizes, depth)
		log.Printf("%v return %v", strings.Repeat(" ", depth), total)
		return total
	default:
		log.Fatal("Unknown item in sequence")
		return 0
	}
}

func countArrangementsIfPresent(sequence []spring, groupSizes []int, depth int) int {
	if len(groupSizes) == 0 {
		log.Printf("%v can't consume present, empty group list", strings.Repeat(" ", depth))
		return 0
	}
	groupSize := groupSizes[0]
	if groupSize > len(sequence) {
		log.Printf("%v can't consume present, groupSize=%v len(sequence)=%v", strings.Repeat(" ", depth), groupSize, len(sequence))
		return 0
	}
	for i := 1; i < groupSize; i += 1 {
		if sequence[i] == Empty {
			log.Printf("%v found empty spot, return 0", strings.Repeat(" ", depth))
			return 0
		}
	}
	if groupSize == len(sequence) {
		if len(groupSizes) == 1 {
			log.Printf("%v end of sequence, return 1", strings.Repeat(" ", depth))
			return 1
		} else {
			log.Printf("%v end of sequence, return 0", strings.Repeat(" ", depth))
			return 0
		}
	}
	if sequence[groupSize] == Present {
		log.Printf("%v group too long, return 0", strings.Repeat(" ", depth))
		return 0
	}
	log.Printf("%v consume %v present 1 absent", strings.Repeat(" ", depth), groupSize)
	return countArrangements(sequence[groupSize+1:], groupSizes[1:], depth+1)
}
