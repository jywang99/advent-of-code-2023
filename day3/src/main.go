package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"unicode"
)

// var file = "example.txt"
var file = "input.txt"

func isSymbol(c rune) bool {
    return !unicode.IsDigit(c) && c != '.'
}

// get indices of symbols in the line
func getIndices(symline *string) []int {
	indices := make([]int, 0)
	for i, c := range *symline {
		if !unicode.IsDigit(c) && c != '.' {
			indices = append(indices, i)
		}
	}
	return indices
}

// iterator-like struct for parsing lines
type lineParser struct {
	str string
	i   int
}

// value, start, end positions of the number
// start: inclusive, end: exclusive
type intRange struct {
	start int
	end   int
	value int
}

func newLineParser(s string) *lineParser {
	return &lineParser{str: s, i: 0}
}

func (lp *lineParser) nextInt() (intRange, bool) {
	// find first digit
	si := -1
	for ; lp.i < len(lp.str); lp.i++ {
		c := rune(lp.str[lp.i])
		if unicode.IsDigit(c) {
			si = lp.i
			break
		}
	}

	// no digit found
	if si == -1 {
		return intRange{}, false
	}

	// find next non-digit
	ei := len(lp.str)
	for ; lp.i < len(lp.str); lp.i++ {
		c := rune(lp.str[lp.i])
		if !unicode.IsDigit(c) {
			ei = lp.i
			break
		}
	}

	// convert to int
	nStr := lp.str[si:ei]
	n, err := strconv.Atoi(nStr)
	if err != nil {
		log.Fatal(err)
	}
	return intRange{si, ei, n}, true
}

func getRowSum(rows *[3]string) int {
	sum := 0

	// valid indices
	indices := getIndices(&rows[0])
	indices = append(indices, getIndices(&rows[2])...)

	// find numbers and do sum
	parser := newLineParser(rows[1])
	for {
		ir, found := parser.nextInt()
		if !found {
			// no more numbers, done processing row
			break
		}

		// check if the number is valid
		// symbols in same row
		if ir.start > 0 && isSymbol(rune(rows[1][ir.start-1])) || (ir.end < len(rows[1]) && isSymbol(rune(rows[1][ir.end]))) {
			sum += ir.value
			continue
		}

		// symbols in above or below rows
		for _, i := range indices {
			if ir.start-1 <= i && ir.end >= i {
				sum += ir.value
				break
			}
		}
	}

	return sum
}

func slideRows(lines *[3]string, nline *string) {
	lines[0] = lines[1]
	lines[1] = lines[2]
	lines[2] = *nline
}

func processLines(scanner *bufio.Scanner)(int, int, error) {
	// prep reading lines
	lines := [3]string{}
	fsum := 0

	// first row
	scanner.Scan()
    if err := scanner.Err(); err != nil {
        return 0, 0, err
    }
	lines[1] = scanner.Text()
	scanner.Scan()
	lines[2] = scanner.Text()
	fsum += getRowSum(&lines)

	// inbetween
	for scanner.Scan() {
        if err := scanner.Err(); err != nil {
            return 0, 0, err
        }
		lines[0] = lines[1]
		lines[1] = lines[2]
		lines[2] = scanner.Text()
		indices := getIndices(&lines[0])
		indices = append(indices, getIndices(&lines[2])...)
		fsum += getRowSum(&lines)
	}

	// last row
	lines[0] = lines[1]
	lines[1] = lines[2]
	lines[2] = ""
	fsum += getRowSum(&lines)

    return fsum, 0, nil
}

func main() {
	// read file
	fpath := filepath.Join("data", file)
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

    fsum, _, err := processLines(bufio.NewScanner(file))
    if err != nil {
        log.Fatal(err)
    }

	log.Printf("Sum: %d", fsum)
}
