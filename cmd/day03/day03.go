package main

import "bufio"
import "fmt"
import "log"
import "os"

func isSymbol(char rune) bool {
	if char == '.' {
		return false
	} else if '0' <= char && char <= '9' {
		return false
	} else {
		return true
	}
}

func isThereSymbol(line string, fromPos int, toPos int) bool {
	if fromPos < 0 {
		fromPos = 0
	}
	if toPos > len(line) {
		toPos = len(line)
	}
	log.Printf("isThereSymbol(\"%v\")", line[fromPos:toPos])
	for _, char := range line[fromPos:toPos] {
		if isSymbol(char) {
			return true
		}
	}
	return false
}

type Parser struct {
	sumOfPartNumbers int
	sumOfGearRatios  int
}

func (parser *Parser) Parse(lines []string) {
	log.Printf("Parse(\n    \"%v\"\n    \"%v\"\n    \"%v\"\n)", lines[0], lines[1], lines[2])

	// Iterate on the current line to find part numbers
	num := 0
	numStart := -1
	for pos, char := range lines[1] + "." {
		if '0' <= char && char <= '9' {
			// Read the number
			if numStart == -1 {
				numStart = pos
				num = 0
			}
			num = num*10 + int(char-'0')
		} else if numStart != -1 {
			log.Printf("Number %v from %v to %v", num, numStart, pos)
			// We read a whole number
			// Look whether there is a symbol somewhere around
			isPart := isThereSymbol(lines[1], numStart-1, pos+1)
			if len(lines[0]) > 0 {
				isPart = isPart || isThereSymbol(lines[0], numStart-1, pos+1)
			}
			if len(lines[2]) > 0 {
				isPart = isPart || isThereSymbol(lines[2], numStart-1, pos+1)
			}

			if isPart {
				log.Print("PART: YES")
				parser.sumOfPartNumbers += num
			} else {
				log.Print("PART: NO")
			}

			numStart = -1
		}
	}

	// Iterate on the current line to find gears
	for gearPos, char := range lines[1] + "." {
		if char != '*' {
			continue
		}
		var neighboringNumbers []int
		for _, line := range lines {
			num := 0
			numStart := -1

			// Iterate on each line to find part numbers
			for pos, char := range line + "." {
				if '0' <= char && char <= '9' {
					// Read the number
					if numStart == -1 {
						numStart = pos
						num = 0
					}
					num = num*10 + int(char-'0')
				} else if numStart != -1 {
					// Add to the list if adjacent
					if numStart-1 <= gearPos && gearPos <= pos {
						neighboringNumbers = append(neighboringNumbers, num)
					}
					numStart = -1
				}
			}
		}

		// If exactly two, multiply
		if len(neighboringNumbers) == 2 {
			log.Printf("found gear pos %v: neighbors %v, %v", gearPos, neighboringNumbers[0], neighboringNumbers[1])
			parser.sumOfGearRatios += neighboringNumbers[0] * neighboringNumbers[1]
		}
	}
}

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day03.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	parser := Parser{
		sumOfPartNumbers: 0,
	}

	// Iterate on lines
	var lines [3]string
	for fileScanner.Scan() {
		lines[0] = lines[1]
		lines[1] = lines[2]
		lines[2] = fileScanner.Text()

		if len(lines[1]) > 0 {
			parser.Parse(lines[0:3])
		}
	}
	lines[0] = lines[1]
	lines[1] = lines[2]
	lines[2] = ""
	parser.Parse(lines[0:3])

	fmt.Println(parser.sumOfPartNumbers)
	fmt.Println(parser.sumOfGearRatios)
}
