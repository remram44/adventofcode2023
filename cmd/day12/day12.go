package main

import "bufio"
import "fmt"
import "log"
import "os"
import "runtime"
//import "strings"
import "sync"
import "sync/atomic"

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

	// Channel for the inputs
	lineInputs := make(chan string)

	// Wait group to notice the end
	var arrangementsWg sync.WaitGroup

	var sumOfArrangementsPart1 int64 = 0
	var sumOfArrangementsPart2 int64 = 0

	// Create the processing threads
	numCPU := runtime.GOMAXPROCS(0)
	log.Printf("Starting %v goroutines", numCPU)
	for i := 0; i < numCPU; i += 1 {
		go func() {
			for {
				line := <-lineInputs

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

				// Find arrangements for part 1
				//log.Print(line)
				arrangements1 := countArrangements(sequence, groupSizes, 0)
				//log.Printf("%v -> %v", line, arrangements1)
				//log.Print()
				atomic.AddInt64(&sumOfArrangementsPart1, int64(arrangements1))

				// Unfold the sequence
				var sequence2 []spring
				for i, element := range sequence {
					if i != 0 {
						sequence2 = append(sequence2, Unknown)
					}
					for j := 0; j < 5; j += 1 {
						sequence2 = append(sequence2, element)
					}
				}

				// Unfold the groups
				var groupSizes2 []int
				for _, element := range groupSizes {
					for j := 0; j < 5; j += 1 {
						groupSizes2 = append(groupSizes2, element)
					}
				}

				// Find arrangements for part 2
				//log.Printf("%v x5", line)
				arrangements2 := countArrangements(sequence2, groupSizes2, 0)
				//log.Printf("%v x5 -> %v", line, arrangements2)
				//log.Print()
				atomic.AddInt64(&sumOfArrangementsPart2, int64(arrangements2))

				arrangementsWg.Done()
			}
		}()
	}

	// Iterate on lines
	lineNum := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineNum += 1
		arrangementsWg.Add(1)
		lineInputs <- line
	}

	arrangementsWg.Wait()

	fmt.Println(sumOfArrangementsPart1)

	fmt.Println(sumOfArrangementsPart2)
}

func countArrangements(sequence []spring, groupSizes []int, depth int) int {
	//log.Printf("%vcountArrangements(%v, %v)", strings.Repeat(" ", depth), sequence, groupSizes)
	if len(sequence) == 0 {
		if len(groupSizes) == 0 {
			//log.Printf("%v done, return 1", strings.Repeat(" ", depth))
			return 1
		} else {
			//log.Printf("%v fail, return 0", strings.Repeat(" ", depth))
			return 0
		}
	}

	switch sequence[0] {
	case Empty:
		//log.Printf("%v empty, consume sequence", strings.Repeat(" ", depth))
		result := countArrangements(sequence[1:], groupSizes, depth+1)
		//log.Printf("%v return %v", strings.Repeat(" ", depth), result)
		return result
	case Present:
		result := countArrangementsIfPresent(sequence, groupSizes, depth)
		//log.Printf("%v return %v", strings.Repeat(" ", depth), result)
		return result
	case Unknown:
		total := 0
		//log.Printf("%v try guess absent", strings.Repeat(" ", depth))
		total += countArrangements(sequence[1:], groupSizes, depth+1)
		//log.Printf("%v try guess present", strings.Repeat(" ", depth))
		total += countArrangementsIfPresent(sequence, groupSizes, depth)
		//log.Printf("%v return %v", strings.Repeat(" ", depth), total)
		return total
	default:
		log.Fatal("Unknown item in sequence")
		return 0
	}
}

func countArrangementsIfPresent(sequence []spring, groupSizes []int, depth int) int {
	if len(groupSizes) == 0 {
		//log.Printf("%v can't consume present, empty group list", strings.Repeat(" ", depth))
		return 0
	}
	groupSize := groupSizes[0]
	if groupSize > len(sequence) {
		//log.Printf("%v can't consume present, groupSize=%v len(sequence)=%v", strings.Repeat(" ", depth), groupSize, len(sequence))
		return 0
	}
	for i := 1; i < groupSize; i += 1 {
		if sequence[i] == Empty {
			//log.Printf("%v found empty spot, return 0", strings.Repeat(" ", depth))
			return 0
		}
	}
	if groupSize == len(sequence) {
		if len(groupSizes) == 1 {
			//log.Printf("%v end of sequence, return 1", strings.Repeat(" ", depth))
			return 1
		} else {
			//log.Printf("%v end of sequence, return 0", strings.Repeat(" ", depth))
			return 0
		}
	}
	if sequence[groupSize] == Present {
		//log.Printf("%v group too long, return 0", strings.Repeat(" ", depth))
		return 0
	}
	//log.Printf("%v consume %v present 1 absent", strings.Repeat(" ", depth), groupSize)
	return countArrangements(sequence[groupSize+1:], groupSizes[1:], depth+1)
}
