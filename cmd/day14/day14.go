package main

import "bufio"
import "fmt"
import "log"
import "os"

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

	// Do the rolling
	rollUp(field)
	log.Print("Result:")
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
