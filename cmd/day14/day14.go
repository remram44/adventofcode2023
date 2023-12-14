package main

import "bufio"
import "fmt"
import "log"
import "os"
import "slices"

type tile int

const (
	Empty tile = iota
	Fixed
	Rolling
)

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day14.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines
	var field [][]tile
	for fileScanner.Scan() {
		line := fileScanner.Text()
		var fieldLine []tile
		for _, c := range line {
			switch c {
			case '.':
				fieldLine = append(fieldLine, Empty)
			case '#':
				fieldLine = append(fieldLine, Fixed)
			case 'O':
				fieldLine = append(fieldLine, Rolling)
			default:
				log.Fatalf("Unknown character %c in field", c)
			}
		}
		field = append(field, fieldLine)
	}

	// Part 1
	{
		var fieldPart1 [][]tile = make([][]tile, len(field))[:0]
		for y := 0; y < len(field); y += 1 {
			fieldPart1 = append(fieldPart1, slices.Clone(field[y]))
		}

		// Do the rolling
		rollUp(fieldPart1)
		log.Print("Result part 1:")
		printField(fieldPart1)

		// Count the load
		load := measureNorthLoad(fieldPart1)

		fmt.Println(load)
	}

	// Part 2
	previousCycles := make(map[string]int)
	target := 1000000000
	for cycle := 0; cycle < target; cycle += 1 {
		key := fieldKey(field)
		prevCycle, ok := previousCycles[key]
		if !ok {
			previousCycles[key] = cycle
		} else {
			// There's a repeat every (cycle-prevCycle) cycles
			log.Printf("repeat %v %v", prevCycle, cycle)
			if target-cycle > cycle-prevCycle {
				// Jump forward
				cycle = target - (target-cycle)%(cycle-prevCycle) - 1
				continue
			}
		}

		for dir := 0; dir < 4; dir += 1 {
			rollUp(field)
			field = rotate(field)
		}

	}

	log.Print("Result part 2:")
	printField(field)

	// Count the load
	load := measureNorthLoad(field)

	fmt.Println(load)
}

func printField(field [][]tile) {
	sizeX := len(field[0])
	sizeY := len(field)

	for y := 0; y < sizeY; y += 1 {
		line := ""
		for x := 0; x < sizeX; x += 1 {
			switch field[y][x] {
			case Empty:
				line += "."
			case Fixed:
				line += "#"
			case Rolling:
				line += "O"
			default:
				line += "?"
			}
		}
		log.Print(line)
	}
}

func rollUp(field [][]tile) {
	sizeX := len(field[0])
	sizeY := len(field)

	rollTo := make([]int, sizeX)

	for y := 0; y < sizeY; y += 1 {
		for x := 0; x < sizeX; x += 1 {
			switch field[y][x] {
			case Fixed:
				rollTo[x] = y + 1
			case Rolling:
				// Roll it up
				field[y][x] = Empty
				field[rollTo[x]][x] = Rolling

				rollTo[x] += 1
			}
		}
	}
}

func measureNorthLoad(field [][]tile) int {
	sizeX := len(field[0])
	sizeY := len(field)

	load := 0

	for y := 0; y < sizeY; y += 1 {
		for x := 0; x < sizeX; x += 1 {
			if field[y][x] == Rolling {
				load += sizeY - y
			}
		}
	}

	return load
}

func rotate(field [][]tile) [][]tile {
	sizeX := len(field[0])
	sizeY := len(field)

	var newField [][]tile

	for x := 0; x < sizeX; x += 1 {
		var newFieldLine []tile
		for y := sizeY - 1; y >= 0; y -= 1 {
			newFieldLine = append(newFieldLine, field[y][x])
		}
		newField = append(newField, newFieldLine)
	}

	return newField
}

func fieldKey(field [][]tile) string {
	sizeX := len(field[0])
	sizeY := len(field)

	key := ""

	for y := 0; y < sizeY; y += 1 {
		for x := 0; x < sizeX; x += 1 {
			switch field[y][x] {
			case Empty:
				key += "."
			case Fixed:
				key += "#"
			case Rolling:
				key += "O"
			default:
				key += "?"
			}
		}
		key += "\n"
	}

	return key
}
