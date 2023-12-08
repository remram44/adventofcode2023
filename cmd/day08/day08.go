package main

import "bufio"
import "fmt"
import "log"
import "os"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day08.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Load the directions
	var directions []int
	if !fileScanner.Scan() {
		log.Fatal("Can't read first line")
	}
	line := fileScanner.Text()
	for _, char := range line {
		switch char {
		case 'L':
			directions = append(directions, 0)
		case 'R':
			directions = append(directions, 1)
		default:
			log.Fatalf("Invalid direction %c", char)
		}
	}

	if !fileScanner.Scan() || fileScanner.Text() != "" {
		log.Fatal("Can't read second line")
	}

	// Iterate on lines to load network
	network := make(map[string][2]string)
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if len(line) != 16 || line[3:7] != " = (" || line[10:12] != ", " || line[15:16] != ")" {
			log.Fatalf("Invalid line: %v", line)
		}

		fromNode := line[0:3]
		toNodeL := line[7:10]
		toNodeR := line[12:15]
		network[fromNode] = [2]string{toNodeL, toNodeR}
	}

	// Follow the directions
	steps := 0
	currentNode := "AAA"
	for currentNode != "ZZZ" {
		direction := directions[steps%len(directions)]
		steps += 1
		currentNode = network[currentNode][direction]
	}

	fmt.Println(steps)
}
