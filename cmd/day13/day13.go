package main

import "bufio"
import "fmt"
import "log"
import "os"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day13.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines
	sum := 0
	var pattern [][]bool
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) == 0 {
			sum += summarizePattern(pattern)
			pattern = nil
			continue
		}
		var patternLine []bool
		for _, c := range line {
			switch c {
			case '#':
				patternLine = append(patternLine, true)
			case '.':
				patternLine = append(patternLine, false)
			default:
				log.Fatalf("Unknown character %c in pattern", c)
			}
		}
		pattern = append(pattern, patternLine)
	}
	sum += summarizePattern(pattern)

	fmt.Println(sum)
}

func summarizePattern(pattern [][]bool) int {
	sizeX := len(pattern[0])
	sizeY := len(pattern)

	// Find a vertical line of reflection
	for symX := 1; symX < sizeX; symX += 1 {
		matches := true
		for i := 0; matches && symX-1-i >= 0 && symX+i < sizeX; i += 1 {
			for y := 0; y < sizeY; y += 1 {
				if pattern[y][symX-1-i] != pattern[y][symX+i] {
					matches = false
					break
				}
			}
		}

		if matches {
			// Found it
			log.Printf("  vertical line of reflection %v", symX)
			return symX
		}
	}

	// Find a horizontal line of reflection
	for symY := 1; symY < sizeY; symY += 1 {
		matches := true
		for i := 0; matches && symY-1-i >= 0 && symY+i < sizeY; i += 1 {
			for x := 0; x < sizeX; x += 1 {
				if pattern[symY-1-i][x] != pattern[symY+i][x] {
					matches = false
					break
				}
			}
		}

		if matches {
			// Found it
			log.Printf("horizontal line of reflection %v", symY)
			return 100 * symY
		}
	}

	log.Fatal("Not reflection!")
	return 0
}
