package main

import "bufio"
import "fmt"
import "log"
import "math"
import "os"
import "strings"

import "github.com/remram44/adventofcode2023"

// time
// distance
// holdButton
// distance =  holdButton * (time - holdButton)
// distance = holdButton*time - holdButton^2
// win if: distanceRecord <= holdButton*time - holdButton^2
// win if: holdButton^2 - holdButton*time + distanceRecord <= 0
// roots:
//   delta = (time^2 - 4*distanceRecord)
//   if delta >= 0:
//	   (time - sqrt(delta))/2
//	   (time + sqrt(delta))/2
// For the example:
//   time=7 distanceRecord=9
//   delta=6.5
//   roots: 1.7, 5.3
//   so win for 2, 3, 4, 5

func countRaceWinOptions(time int, distanceRecord int) int {
	timef := float64(time)
	delta := timef*timef - 4.0*float64(distanceRecord)
	if delta < 0 {
		log.Printf("Can't win race: time=%v distanceRecord=%v", time, distanceRecord)
		return 0
	}

	root1 := (timef - math.Sqrt(delta)) / 2.0
	root2 := (timef + math.Sqrt(delta)) / 2.0

	return int(math.Floor(root2)) - int(math.Ceil(root1)) + 1
}

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day06.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Read times
	var times []int
	if !fileScanner.Scan() {
		log.Fatal("Missing line")
	}
	line := fileScanner.Text()
	if !strings.HasPrefix(line, "Time: ") {
		log.Fatal("Invalid times")
	}
	pos := 6
	for pos < len(line) {
		pos = aoc.ReadSpaces(line, pos)
		var num int
		pos, num = aoc.ReadNumber(line, pos)
		times = append(times, num)
	}

	// Read distance records
	var distanceRecords []int
	if !fileScanner.Scan() {
		log.Fatal("Missing line")
	}
	line = fileScanner.Text()
	if !strings.HasPrefix(line, "Distance: ") {
		log.Fatal("Invalid distances")
	}
	pos = 10
	for pos < len(line) {
		pos = aoc.ReadSpaces(line, pos)
		var num int
		pos, num = aoc.ReadNumber(line, pos)
		distanceRecords = append(distanceRecords, num)
	}

	log.Printf("times=%v distances=%v", times, distanceRecords)

	result := 1
	for i := 0; i < len(times); i += 1 {
		result *= countRaceWinOptions(times[i], distanceRecords[i])
	}

	fmt.Println(result)
}
