package main

import "bufio"
import "fmt"
import "log"
import "os"
import "slices"
import "sort"

import "github.com/remram44/adventofcode2023"

type handType int

const (
	highCard handType = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

func (h handType) String() string {
	switch h {
	case highCard:
		return "high card"
	case onePair:
		return "one pair"
	case twoPair:
		return "two pair"
	case threeOfAKind:
		return "three of a kind"
	case fullHouse:
		return "full house"
	case fourOfAKind:
		return "four of a kind"
	case fiveOfAKind:
		return "five of a kind"
	default:
		return "unknown"
	}
}

type hand struct {
	handType handType
	hand     string
}

type input struct {
	hand string
	bid  int
}

func newHand(handStr string, part int) hand {
	// Count each card
	cards := make(map[rune]int)
	for _, card := range handStr {
		cards[card] += 1
	}

	// For part 2, jokers are special
	jokers := 0
	if part == 2 {
		jokers = cards['J']
		delete(cards, 'J')
	}

	// Make a sorted array of cards by count
	var counts []int
	for _, count := range cards {
		counts = append(counts, count)
	}
	slices.SortFunc(counts, func(a int, b int) int { return b - a })

	var handType handType
	switch {
	case jokers == 5:
		handType = fiveOfAKind
	case counts[0]+jokers == 5:
		handType = fiveOfAKind
	case counts[0]+jokers == 4:
		handType = fourOfAKind
	case slices.Equal(counts, []int{3, 2}):
		handType = fullHouse
	case slices.Equal(counts, []int{2, 2}) && jokers == 1:
		handType = fullHouse
	case slices.Equal(counts, []int{3, 1, 1}):
		handType = threeOfAKind
	case slices.Equal(counts, []int{2, 1, 1}) && jokers == 1:
		handType = threeOfAKind
	case slices.Equal(counts, []int{1, 1, 1}) && jokers == 2:
		handType = threeOfAKind
	case slices.Equal(counts, []int{2, 2, 1}):
		handType = twoPair
	case slices.Equal(counts, []int{2, 1, 1, 1}):
		handType = onePair
	case slices.Equal(counts, []int{1, 1, 1, 1}) && jokers == 1:
		handType = onePair
	case slices.Equal(counts, []int{1, 1, 1, 1, 1}):
		handType = highCard
	default:
		log.Fatalf("Can't figure out hand type: %v %v jokers=%v", handStr, counts, jokers)
	}

	return hand{
		handType: handType,
		hand:     handStr,
	}
}

type handBid struct {
	hand hand
	bid  int
}

func faceRank(card byte, part int) int {
	switch card {
	case 'A':
		return 14
	case 'K':
		return 13
	case 'Q':
		return 12
	case 'J':
		if part == 1 {
			return 11
		} else {
			// Weakest card in part 2
			return 1
		}
	case 'T':
		return 10
	default:
		return int(card) - '0'
	}
}

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day07.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Iterate on lines to read bids
	var inputs []input
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) < 6 {
			log.Fatalf("Line too short: %v", line)
		}
		hand := line[0:5]
		if line[5] != ' ' {
			log.Fatal("Missing separator")
		}
		pos, bid := aoc.ReadNumber(line, 6)
		if pos != len(line) {
			log.Fatal("Extra content on line")
		}

		inputs = append(inputs, input{hand, bid})
	}

	for part := 1; part <= 2; part += 1 {
		// Parse hands, make bids array
		var bids []handBid
		for _, input := range inputs {
			bids = append(bids, handBid{newHand(input.hand, part), input.bid})
		}

		// Sort bids
		sort.Slice(bids, func(a int, b int) bool {
			// Compare hand types
			if bids[a].hand.handType < bids[b].hand.handType {
				return true
			} else if bids[a].hand.handType > bids[b].hand.handType {
				return false
			}

			// Compare the hands themselves
			for i := 0; i < 5; i += 1 {
				faceA := faceRank(bids[a].hand.hand[i], part)
				faceB := faceRank(bids[b].hand.hand[i], part)
				if faceA < faceB {
					return true
				} else if faceA > faceB {
					return false
				}
			}
			return false
		})

		// Compute the total score
		score := 0
		for rank, bid := range bids {
			fmt.Printf("%v %v\n", bid.hand.hand, bid.hand.handType)
			score += (rank + 1) * bid.bid
		}

		fmt.Println(score)
	}
}
