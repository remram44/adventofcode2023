package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strings"

type dirFlag int

const (
	left  dirFlag = 1
	right         = 2
	up            = 4
	down          = 8
)

func parseTile(c rune) dirFlag {
	switch c {
	case '|':
		return up | down
	case '-':
		return left | right
	case 'L':
		return up | right
	case 'J':
		return up | left
	case '7':
		return left | down
	case 'F':
		return right | down
	case '.':
		return 0
	case 'S':
		return 0
	default:
		log.Fatalf("Unknown tile %v", c)
		return 0
	}
}

type direction int

func reverseDir(dir direction) direction {
	return (dir + 2) % 4
}

func dirToFlag(dir direction) dirFlag {
	switch dir {
	case 0:
		return right
	case 1:
		return down
	case 2:
		return left
	case 3:
		return up
	default:
		log.Fatalf("Unknown direction %v", dir)
		return 0
	}
}

func flagToDir(flag dirFlag) direction {
	switch flag {
	case right:
		return 0
	case down:
		return 1
	case left:
		return 2
	case up:
		return 3
	default:
		log.Fatalf("Flag is not a direction: %v", flag)
		return 0
	}
}

type Position struct {
	x int
	y int
}

func move(pos Position, dir direction) Position {
	switch dir {
	case 0: // right
		return Position{x: pos.x + 1, y: pos.y}
	case 1: // down
		return Position{x: pos.x, y: pos.y + 1}
	case 2: // left
		return Position{x: pos.x - 1, y: pos.y}
	case 3: // up
		return Position{x: pos.x, y: pos.y - 1}
	default:
		log.Fatalf("Invalid dir %v", dir)
		return Position{x: 0, y: 0}
	}
}

func (dir direction) String() string {
	switch dir {
	case 0:
		return "right"
	case 1:
		return "down"
	case 2:
		return "left"
	case 3:
		return "up"
	default:
		return "unknown"
	}
}

func (flag dirFlag) String() string {
	var result []string
	switch {
	case flag&up != 0:
		result = append(result, "up")
	case flag&left != 0:
		result = append(result, "left")
	case flag&right != 0:
		result = append(result, "right")
	case flag&down != 0:
		result = append(result, "down")
	}
	if len(result) > 0 {
		return "dirFlag:" + strings.Join(result, "|")
	} else {
		return "dirFlag:none"
	}
}

func (pos Position) String() string {
	return fmt.Sprintf("%v;%v", pos.x, pos.y)
}

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day10.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines, load the entire field
	field := make(map[Position]dirFlag)
	y := 0
	var startingPos Position
	sizeX := 0
	sizeY := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		for x, c := range line {
			field[Position{x: x, y: y}] = parseTile(c)
			if c == 'S' {
				startingPos = Position{x: x, y: y}
			}
			sizeX = x + 1
		}
		y += 1
		sizeY = y
	}

	log.Printf("Read field %vx%v", sizeX, sizeY)

	// Start at S and go each direction
	var startDir direction
	var distancesPerDir []map[Position]int
	for startDir = 0; startDir < 4; startDir += 1 {
		pos := startingPos
		dir := startDir

		//log.Printf("Starting at %v dir=%v", pos, dir)

		distance := 1
		distancesPerDir = append(distancesPerDir, make(map[Position]int))
		for {
			//log.Printf("%v dir=%v", pos, dir)

			// Can we move?
			nextPos := move(pos, dir)
			nextPosFlag := field[nextPos]
			flagDirWeComeFrom := dirToFlag(reverseDir(dir))
			if flagDirWeComeFrom&nextPosFlag == 0 {
				//log.Printf("Can't go %v to %v", dir, field[nextPos])
				break
			}

			// Move
			pos = nextPos

			// Record
			distancesPerDir[startDir][pos] = distance
			distance += 1

			// Find new direction
			nextPosFlag = nextPosFlag & ^flagDirWeComeFrom
			dir = flagToDir(nextPosFlag)
		}
	}

	// Now the result is the maximum tile that is in two maps
	maxMinDistance := 0
	var maxLoop map[Position]int
	for y := 0; y < sizeY; y += 1 {
		for x := 0; x < sizeX; x += 1 {
			count := 0
			minDistance := sizeX * sizeY
			var dir direction
			var minDirection direction
			for dir = 0; dir < 4; dir += 1 {
				distance, present := distancesPerDir[dir][Position{x: x, y: y}]
				if present {
					count += 1
					if distance < minDistance {
						minDistance = distance
						minDirection = dir
					}
				}
			}

			if count >= 2 && minDistance > maxMinDistance {
				maxMinDistance = minDistance
				maxLoop = distancesPerDir[minDirection]
			}
		}
	}

	fmt.Println(maxMinDistance)

	// For part 2, we build an array 3 times as big as the field
	// plus 1 on each edge
	// Example:
	// F- turns into:
	// . ... ... .
	//
	// . ... ... .
	// . .xx xxx .
	// . .x. ... .
	//
	// . ... ... .
	floodableField := make([]bool, (sizeX*3+2)*(sizeY*3+2))
	set := func(xt, yt, xr, yr int) {
		x := 2 + xt*3 + xr
		y := 2 + yt*3 + yr
		floodableField[y*(sizeX*3+2)+x] = true
	}
	get := func(xt, yt, xr, yr int) bool {
		x := 2 + xt*3 + xr
		y := 2 + yt*3 + yr
		return floodableField[y*(sizeX*3+2)+x]
	}
	for y := 0; y < sizeY; y += 1 {
		for x := 0; x < sizeX; x += 1 {
			pos := Position{x: x, y: y}
			_, present := maxLoop[pos]
			if !present {
				continue
			}
			tile := field[pos]
			if tile&up != 0 {
				set(x, y, 0, 0)
				set(x, y, 0, -1)
			}
			if tile&left != 0 {
				set(x, y, 0, 0)
				set(x, y, -1, 0)
			}
			if tile&right != 0 {
				set(x, y, 0, 0)
				set(x, y, 1, 0)
			}
			if tile&down != 0 {
				set(x, y, 0, 0)
				set(x, y, 0, 1)
			}
		}
	}

	// Add the start
	for y := -1; y <= 1; y += 1 {
		for x := -1; x <= 1; x += 1 {
			set(startingPos.x, startingPos.y, x, y)
		}
	}

	// Print the floodable field
	for y := 0; y < sizeY*3+2; y += 1 {
		line := ""
		for x := 0; x < sizeX*3+2; x += 1 {
			if floodableField[y*(sizeX*3+2)+x] {
				line += "."
			} else {
				line += " "
			}
		}
		log.Print(line, "<")
	}

	// Now flood
	flood(floodableField, sizeX*3+2, sizeY*3+2)

	// Print the area
	for y := 0; y < sizeY*3+2; y += 1 {
		line := ""
		for x := 0; x < sizeX*3+2; x += 1 {
			if floodableField[y*(sizeX*3+2)+x] {
				line += "."
			} else if x%3 == 2 && y%3 == 2 {
				line += "I"
			} else {
				line += " "
			}
		}
		log.Print(line, "<")
	}

	// Count the area
	area := 0
	for y := 0; y < sizeY; y += 1 {
		for x := 0; x < sizeX; x += 1 {
			if !get(x, y, 0, 0) {
				area += 1
			}
		}
	}

	fmt.Println(area)
}

func flood(field []bool, sizeX int, sizeY int) {
	openList := make(map[Position]interface{})
	openList[Position{x: 0, y: 0}] = nil
	for len(openList) > 0 {
		for pos := range openList {
			delete(openList, pos)
			if !field[pos.y*sizeX+pos.x] {
				field[pos.y*sizeX+pos.x] = true
				for y := -1; y <= 1; y += 1 {
					for x := -1; x <= 1; x += 1 {
						neighbor := Position{x: pos.x + x, y: pos.y + y}
						if neighbor.x >= 0 && neighbor.x < sizeX && neighbor.y >= 0 && neighbor.y < sizeY {
							openList[neighbor] = nil
						}
					}
				}
			}
		}
	}
}
