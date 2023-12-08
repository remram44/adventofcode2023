package main

import "bufio"
import "fmt"
import "log"
import "os"

func runThroughNetwork(directions []int, network map[string][2]string, start string, singleZ bool) int {
	steps := 0
	currentNode := start
	for {
		if singleZ {
			if currentNode[2] == 'Z' {
				break
			}
		} else {
			if currentNode == "ZZZ" {
				break
			}
		}

		direction := directions[steps%len(directions)]
		steps += 1
		currentNode = network[currentNode][direction]
	}

	return steps
}

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

	// Follow the directions in part 1
	fmt.Println(runThroughNetwork(directions, network, "AAA", false))

	// Find starting nodes for part 2, count steps to get to exit, in parallel
	stepsChannel := make(chan int)
	numParallel := 0
	for node := range network {
		if node[2] == 'A' {
			go func(node string) {
				stepsChannel <- runThroughNetwork(directions, network, node, true)
			}(node)
			numParallel += 1
		}
	}

	// Collect into array
	stepsArray := make([]int, numParallel)[0:0]
	for i := 0; i < numParallel; i += 1 {
		stepsArray = append(stepsArray, <-stepsChannel)
	}

	// Part 2 answer is the least common multiple of those numbers
	fmt.Println(leastCommonMultiple(stepsArray))
}

func leastCommonMultiple(numbers []int) int {
	result := numbers[0]

	for i := 1; i < len(numbers); i += 1 {
		result = result * numbers[i] / greatestCommonDivisor(result, numbers[i])
	}

	return result
}

// Euclid's
func greatestCommonDivisor(a int, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
