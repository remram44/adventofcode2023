package main

import "bufio"
import "fmt"
import "log"
import "os"
import "slices"
import "strconv"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day17.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines
	var grid [][]int
	for fileScanner.Scan() {
		line := fileScanner.Text()
		var row []int
		for _, c := range line {
			loss := int(c - '0')
			row = append(row, loss)
		}
		grid = append(grid, row)
	}

	fmt.Println(pathFind(grid))
}

func pathFind(grid [][]int) int {
	width := len(grid[0])
	height := len(grid)

	type pos struct {
		x int
		y int
	}
	type configuration struct {
		x         int
		y         int
		dx        int
		dy        int
		straights int
		totalLoss int
		path      []pos
	}
	openList := []configuration{
		{
			x:         0,
			y:         0,
			dx:        0,
			dy:        0,
			straights: 0,
			totalLoss: 0,
			path:      []pos{{x: 0, y: 0}},
		},
	}
	closedList := make(map[string]struct{})
	for len(openList) > 0 {
		// Pop last element
		config := openList[len(openList)-1]
		openList = openList[0 : len(openList)-1]
		log.Printf("open configs: %v min loss: %v", len(openList), config.totalLoss)

		key := fmt.Sprintf("%v-%v-%v-%v-%v", config.x, config.y, config.dx, config.dy, config.straights)
		_, closed := closedList[key]
		if closed {
			continue
		}
		closedList[key] = struct{}{}

		// Reached target
		if config.x == width-1 && config.y == height-1 {
			// Copy grid
			var newGrid [][]int
			for y := 0; y < height; y += 1 {
				newGrid = append(newGrid, slices.Clone(grid[y]))
			}

			// Mark path
			for _, pos := range config.path {
				newGrid[pos.y][pos.x] = -1
			}

			// Print
			for y := 0; y < height; y += 1 {
				line := ""
				for x := 0; x < height; x += 1 {
					if newGrid[y][x] == -1 {
						line += "#"
					} else {
						line += strconv.Itoa(newGrid[y][x])
					}
				}
				log.Print(line)
			}
			return config.totalLoss
		}

		// Find possibilities
		for _, d := range []struct {
			x int
			y int
		}{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
			// Can't go opposite way
			if d.x == -config.dx && d.y == -config.dy {
				continue
			}

			// Can't go out of bounds
			if config.x+d.x < 0 || config.x+d.x >= width || config.y+d.y < 0 || config.y+d.y >= height {
				continue
			}

			straights := 1
			if d.x == config.dx && d.y == config.dy {
				straights = config.straights + 1

				if straights > 3 {
					// Can't go more than 3 blocks straight
					continue
				}
			}

			x := config.x + d.x
			y := config.y + d.y
			newConfig := configuration{
				straights: straights,
				x:         x,
				y:         y,
				dx:        d.x,
				dy:        d.y,
				totalLoss: config.totalLoss + grid[y][x],
				path:      append(slices.Clone(config.path), pos{x: x, y: y}),
			}

			// Find insertion point such that openList is in decreasing order of totalLoss
			idx, _ := slices.BinarySearchFunc(
				openList,
				newConfig,
				func(a configuration, b configuration) int {
					return b.totalLoss - a.totalLoss
				},
			)
			openList = slices.Insert(openList, idx, newConfig)
		}
	}

	return -1
}
