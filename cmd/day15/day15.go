package main

import "bufio"
import "fmt"
import "log"
import "os"
import "slices"
import "strconv"

func main() {
	// Open the input file
	readFile, err := os.Open("inputs/day15.txt")
	if err != nil {
		log.Fatalf("Can't open input: %v", err)
	}
	defer readFile.Close()

	// Create a line-by-line scanner
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Read the one line
	if !fileScanner.Scan() {
		log.Fatal("No line in file")
	}
	line := fileScanner.Text()
	if fileScanner.Scan() {
		log.Fatal("Multiple lines in file")
	}

	// Part 1
	{
		// Iterate on comma-separated steps
		start := 0
		sumOfHashes := 0
		for pos := 0; pos < len(line); pos += 1 {
			if line[pos] == ',' {
				sumOfHashes += int(hash(line[start:pos]))
				start = pos + 1
			}
		}
		if start < len(line) {
			sumOfHashes += int(hash(line[start:]))
		}

		fmt.Println(sumOfHashes)
	}

	// Part 2
	{
		type lens struct {
			label       string
			focalLength int
		}

		// Create the 256 boxes
		boxes := make([][]lens, 256)

		equalSign := func(label string, focalLength int) {
			h := int(hash(label))
			idx := slices.IndexFunc(
				boxes[h],
				func(l lens) bool {
					return l.label == label
				},
			)
			if idx != -1 {
				boxes[h][idx].focalLength = focalLength
			} else {
				boxes[h] = append(
					boxes[h],
					lens{
						label:       label,
						focalLength: focalLength,
					},
				)
			}
		}

		dash := func(label string) {
			h := int(hash(label))
			idx := slices.IndexFunc(
				boxes[h],
				func(l lens) bool {
					return l.label == label
				},
			)
			if idx != -1 {
				boxes[h] = slices.Delete(boxes[h], idx, idx+1)
			}
		}

		printBoxes := func(step string) {
			log.Printf("After \"%v\":", step)
			for i, box := range boxes {
				if len(box) > 0 {
					content := ""
					for _, l := range box {
						content += fmt.Sprintf(" [%v %v]", l.label, l.focalLength)
					}
					log.Printf("Box %v:%v", i, content)
				}
			}
		}

		handle := func(step string) {
			for i := 0; i < len(step); i += 1 {
				if step[i] == '=' {
					focalLength, err := strconv.Atoi(step[i+1:])
					if err != nil {
						log.Fatalf("Invalid focal length %#v", step[i+1:])
					}
					equalSign(step[:i], focalLength)
					printBoxes(step)
					return
				} else if step[i] == '-' {
					dash(step[:i])
					printBoxes(step)
					return
				}
			}
			log.Fatalf("Invalid step %#v", step)
		}

		// Iterate on comma-separated steps
		start := 0
		for pos := 0; pos < len(line); pos += 1 {
			if line[pos] == ',' {
				handle(line[start:pos])
				start = pos + 1
			}
		}
		if start < len(line) {
			handle(line[start:])
		}

		// Compute result
		sumOfFocusingPower := 0
		for boxNum, box := range boxes {
			for lensNum, lens := range box {
				sumOfFocusingPower += (1 + boxNum) * (1 + lensNum) * lens.focalLength
			}
		}
		fmt.Println(sumOfFocusingPower)
	}
}

func hash(s string) byte {
	var result byte = 0
	for i := 0; i < len(s); i += 1 {
		result += s[i]
		result *= 17
	}
	return result
}
