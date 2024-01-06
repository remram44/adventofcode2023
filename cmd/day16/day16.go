package main

import "bufio"
import "fmt"
import "log"
import "os"

type tile byte

const (
	empty tile = iota
	splitHorizontal
	splitVertical
	mirrorSlash
	mirrorBackslash
)

type beamDirection int

const (
	right = 0x01
	left  = 0x02
	down  = 0x04
	up    = 0x08
)

type beam struct {
	x         int
	y         int
	direction beamDirection
}

func advance(b beam) beam {
	switch b.direction {
	case right:
		b.x += 1
	case left:
		b.x -= 1
	case down:
		b.y += 1
	case up:
		b.y -= 1
	default:
		log.Fatal("Invalid beam")
	}
	return b
}

func turn(dir beamDirection) beamDirection {
	switch dir {
	case right:
		return down
	case left:
		return up
	case down:
		return left
	case up:
		return right
	default:
		log.Fatal("Invalid direction to turn")
		return 0
	}
}

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day16.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines to read the grid
	var grid [][]tile
	for fileScanner.Scan() {
		line := fileScanner.Text()
		var row []tile
		for _, c := range line {
			var t tile
			switch c {
			case '.':
				t = empty
			case '-':
				t = splitHorizontal
			case '|':
				t = splitVertical
			case '/':
				t = mirrorSlash
			case '\\':
				t = mirrorBackslash
			default:
				log.Fatalf("Unknown tile %#v", c)
			}
			row = append(row, t)
		}
		grid = append(grid, row)
	}

	fmt.Println(energize(
		grid,
		beam{
			x:         0,
			y:         0,
			direction: right,
		},
	))
}

func energize(grid [][]tile, startBeam beam) int {
	width := len(grid[0])
	height := len(grid)

	// Create another array to store energy beams
	var energized = make([][]beamDirection, height)
	for i := 0; i < len(grid); i += 1 {
		energized[i] = make([]beamDirection, width)
	}

	// Flood energy through
	openList := []beam{startBeam}
	appendIfInRange := func(list *[]beam, beam beam) {
		if beam.y >= 0 && beam.y < height && beam.x >= 0 && beam.x < width {
			*list = append(*list, beam)
		}
	}
	for len(openList) > 0 {
		currentList := openList
		openList = nil
		for _, beam := range currentList {
			// Already energized? Done
			if energized[beam.y][beam.x]&beam.direction != 0 {
				continue
			}

			// Energize
			energized[beam.y][beam.x] |= beam.direction

			// Propagate
			switch grid[beam.y][beam.x] {
			case empty:
				// Straight
				appendIfInRange(&openList, advance(beam))
			case splitHorizontal:
				if beam.direction == left || beam.direction == right {
					// Hit the pointy side, go straight
					appendIfInRange(&openList, advance(beam))
				} else {
					// Hit the flat side, split
					newBeam := beam
					newBeam.direction = turn(newBeam.direction)
					appendIfInRange(&openList, advance(newBeam))
					newBeam.direction = turn(turn(newBeam.direction))
					appendIfInRange(&openList, advance(newBeam))
				}
			case splitVertical:
				if beam.direction == up || beam.direction == down {
					// Hit the pointy side, go straight
					appendIfInRange(&openList, advance(beam))
				} else {
					// Hit the flat side, split
					newBeam := beam
					newBeam.direction = turn(newBeam.direction)
					appendIfInRange(&openList, advance(newBeam))
					newBeam.direction = turn(turn(newBeam.direction))
					appendIfInRange(&openList, advance(newBeam))
				}
			case mirrorSlash:
				switch beam.direction {
				case right:
					beam.direction = up
				case left:
					beam.direction = down
				case up:
					beam.direction = right
				case down:
					beam.direction = left
				}
				appendIfInRange(&openList, advance(beam))
			case mirrorBackslash:
				switch beam.direction {
				case right:
					beam.direction = down
				case left:
					beam.direction = up
				case up:
					beam.direction = left
				case down:
					beam.direction = right
				}
				appendIfInRange(&openList, advance(beam))
			}
		}
	}

	// Print beam like in example
	for y := 0; y < height; y += 1 {
		line := ""
		for x := 0; x < width; x += 1 {
			var c rune = '?'
			switch grid[y][x] {
			case splitHorizontal:
				c = '-'
			case splitVertical:
				c = '|'
			case mirrorSlash:
				c = '/'
			case mirrorBackslash:
				c = '\\'
			case empty:
				switch energized[y][x] {
				case right:
					c = '>'
				case left:
					c = '<'
				case up:
					c = '^'
				case down:
					c = 'v'
				case 0:
					c = '.'
				default:
					count := 0
					if energized[y][x]&right != 0 {
						count += 1
					}
					if energized[y][x]&left != 0 {
						count += 1
					}
					if energized[y][x]&up != 0 {
						count += 1
					}
					if energized[y][x]&down != 0 {
						count += 1
					}
					c = rune('0' + count)
				}
			}
			line += string(c)
		}
		log.Print(line)
	}

	// Count energized tiles
	numEnergized := 0
	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			if energized[y][x] != 0 {
				numEnergized += 1
			}
		}
	}
	return numEnergized
}
