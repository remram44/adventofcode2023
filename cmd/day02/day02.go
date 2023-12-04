package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strings"

import "github.com/remram44/adventofcode2023"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day02.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	sumOfIDs := 0
	sumOfPowers := 0

	// Iterate on lines
	for fileScanner.Scan() {
		gameID, isPossible, power := checkGame(fileScanner.Text())
		if isPossible {
			sumOfIDs += gameID
		}
		sumOfPowers += power
	}

	fmt.Println(sumOfIDs)
	fmt.Println(sumOfPowers)
}

func checkGame(line string) (int, bool, int) {
	// Check prefix
	if !strings.HasPrefix(line, "Game ") {
		log.Fatalf("Invalid line: %v", line)
	}

	// Read game ID
	pos, gameID := aoc.ReadNumber(line, 5)

	if line[pos] != ':' || line[pos+1] != ' ' {
		log.Fatalf("Missing colon at %v: %v", pos, line)
	}
	pos += 2

	minRed, minGreen, minBlue := 0, 0, 0

	// Read revealed cubes
	isPossible := true
	for {
		var red, green, blue int
		pos, red, green, blue = readRevealedCubes(line, pos)

		if red > minRed {
			minRed = red
		}
		if green > minGreen {
			minGreen = green
		}
		if blue > minBlue {
			minBlue = blue
		}

		// Check
		isPossible = isPossible && red <= 12 && green <= 13 && blue <= 14

		if pos >= len(line) {
			break
		} else {
			if pos+1 >= len(line) || line[pos] != ';' || line[pos+1] != ' ' {
				log.Fatalf("Missing semicolon at %v: %v", pos, line)
			}
			pos += 2
		}
	}

	power := minRed * minGreen * minBlue

	return gameID, isPossible, power
}

func readRevealedCubes(line string, pos int) (int, int, int, int) {
	red, green, blue := 0, 0, 0

	for {
		// Read the number of cubes
		var num int
		pos, num = aoc.ReadNumber(line, pos)
		if line[pos] != ' ' {
			log.Fatalf("Missing space at %v: %v", pos, line)
		}
		pos += 1

		// Read the color of cubes
		if strings.HasPrefix(line[pos:], "red") {
			red += num
			pos += 3
		} else if strings.HasPrefix(line[pos:], "green") {
			green += num
			pos += 5
		} else if strings.HasPrefix(line[pos:], "blue") {
			blue += num
			pos += 4
		} else {
			log.Fatalf("Unknown color at %v: %v", pos, line)
		}

		// Continue?
		if pos >= len(line) || line[pos] == ';' {
			break
		} else if line[pos] == ',' {
			if line[pos+1] != ' ' {
				log.Fatalf("Unexpected end-of-reveal character at %v: %v", pos, line)
			}
			pos += 2
		} else {
			log.Fatalf("Unexpected end-of-reveal character at %v: %v", pos, line)
		}
	}

	return pos, red, green, blue
}
