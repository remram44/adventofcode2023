package main

import "bufio"
import "fmt"
import "log"
import "os"

type Position struct {
	x int
	y int
}

func absInt(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day11.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines, read initial universe
	var initialUniverse [][]bool
	for fileScanner.Scan() {
		line := fileScanner.Text()
		row := make([]bool, len(line))
		for i, c := range line {
			switch c {
			case '#':
				row[i] = true
			case '.':
				row[i] = false
			default:
				log.Fatalf("Unknown cosmic object %c", c)
			}
		}
		initialUniverse = append(initialUniverse, row)
	}

	sizeX := len(initialUniverse[0])
	sizeY := len(initialUniverse)
	log.Printf("initialUniverse %vx%v", sizeX, sizeY)

	doWithExpansion(initialUniverse, sizeX, sizeY, 2)
	doWithExpansion(initialUniverse, sizeX, sizeY, 1000000)
}

func doWithExpansion(initialUniverse [][]bool, sizeX int, sizeY int, expansion int) {
	// Find mapping of initial Y pos to expanded pos
	var mapY []int
	dest := 0
	for y := 0; y < sizeY; y += 1 {
		rowIsEmpty := true
		for x := 0; x < sizeX; x += 1 {
			if initialUniverse[y][x] {
				rowIsEmpty = false
				break
			}
		}
		mapY = append(mapY, dest)
		if rowIsEmpty {
			dest += expansion
		} else {
			dest += 1
		}
	}

	// Find mapping of initial X pos to expanded pos
	var mapX []int
	dest = 0
	for x := 0; x < sizeX; x += 1 {
		columnIsEmpty := true
		for y := 0; y < sizeY; y += 1 {
			if initialUniverse[y][x] {
				columnIsEmpty = false
				break
			}
		}
		mapX = append(mapX, dest)
		if columnIsEmpty {
			dest += expansion
		} else {
			dest += 1
		}
	}

	// Build list of galaxies, mapping to expanded pos
	var galaxies []Position
	for y := 0; y < sizeY; y += 1 {
		for x := 0; x < sizeX; x += 1 {
			if initialUniverse[y][x] {
				galaxies = append(galaxies, Position{x: mapX[x], y: mapY[y]})
			}
		}
	}

	// For each pair, find shortest path
	sumShortestDistances := 0
	for i := 0; i < len(galaxies); i += 1 {
		for j := i + 1; j < len(galaxies); j += 1 {
			if i == j {
				continue
			}
			distance := absInt(galaxies[i].x-galaxies[j].x) + absInt(galaxies[i].y-galaxies[j].y)
			sumShortestDistances += distance
		}
	}

	fmt.Println(sumShortestDistances)
}
