package aoc

import "log"

func ReadSpaces(line string, pos int) int {
	for {
		if pos >= len(line) {
			log.Fatal("Reached end of line while reading spaces")
		}
		if line[pos] != ' ' {
			break
		}
		pos += 1
	}
	return pos
}

func ReadNumber(line string, pos int) (int, int) {
	num := 0
	if line[pos] < '0' || line[pos] > '9' {
		log.Fatalf("Not a number at pos %v", pos)
	}
	for pos < len(line) {
		char := line[pos]
		if '0' <= char && char <= '9' {
			num = num*10 + int(char-'0')
		} else {
			break
		}
		pos += 1
	}
	return pos, num
}
