package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strings"

import "github.com/remram44/adventofcode2023"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day04.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines
	winningsPart1 := 0
	cardDuplicates := make(map[int]int)
	totalCardsPart2 := 0
	cardID := 0
	for fileScanner.Scan() {
		cardID += 1

		log.Printf("card ID %v (duplicates: %v)", cardID, cardDuplicates[cardID])
		matchingNumbers := countMatchingNumbers(fileScanner.Text(), cardID)
		log.Printf("matchingNumbers = %v", matchingNumbers)

		// Score for part 1
		winningsPart1 += scoreForMatches(matchingNumbers)

		// Duplicate cards for part 2
		log.Printf("Duplicate next %v cards %v times\n\n", matchingNumbers, 1+cardDuplicates[cardID])
		for i := 1; i <= matchingNumbers; i += 1 {
			cardDuplicates[cardID+i] += 1 + cardDuplicates[cardID]
		}
		totalCardsPart2 += 1 + cardDuplicates[cardID]
	}

	fmt.Println(winningsPart1)
	fmt.Println(totalCardsPart2)
}

func scoreForMatches(matchingNumbers int) int {
	if matchingNumbers == 0 {
		return 0
	} else {
		return 1 << (matchingNumbers - 1)
	}
}

func countMatchingNumbers(line string, expectedCardID int) int {
	log.Printf("countWinnings(\"%v\")", line)

	if !strings.HasPrefix(line, "Card ") {
		log.Fatal("Missing Card prefix")
	}
	pos := 5

	pos = aoc.ReadSpaces(line, pos)

	var cardID int
	pos, cardID = aoc.ReadNumber(line, pos)
	if cardID != expectedCardID {
		log.Fatalf("Unexpected cardID: %v != %v", cardID, expectedCardID)
	}

	if pos >= len(line) || line[pos] != ':' {
		log.Fatal("Missing colon")
	}
	pos += 1

	// Read winning numbers
	winningNumbers := make(map[int]bool)
	for {
		pos = aoc.ReadSpaces(line, pos)
		if line[pos] == '|' {
			pos += 1
			break
		}
		var num int
		pos, num = aoc.ReadNumber(line, pos)
		winningNumbers[num] = true
		log.Printf("  winning number: %v", num)
	}

	// Read our numbers
	matchingNumbers := 0
	for pos < len(line) {
		pos = aoc.ReadSpaces(line, pos)
		var num int
		pos, num = aoc.ReadNumber(line, pos)
		if winningNumbers[num] {
			matchingNumbers += 1
			log.Printf("  matching number: %v", num)
		} else {
			log.Printf("  unmatched number: %v", num)
		}
	}

	return matchingNumbers
}
