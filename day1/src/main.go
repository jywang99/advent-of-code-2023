package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// var file = "example.txt"
var file = "input.txt"
// var file = "example2.txt"

var nMap = map[string]int{
	"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "zero": 0,
}
var maxLen = func() int {
	max := 0
	for k := range nMap {
		if len(k) > max {
			max = len(k)
		}
	}
	return max
}()

// isPossibleNumber checks if a string is a possible number
// return values: (isNumber, isPossible, number)
func isPossibleNumber(s string, prefix bool) (bool, bool, int) {
	// check if string is a number
	if len(s) == 1 {
		c := rune(s[0])
		if unicode.IsDigit(c) {
			return true, true, int(c - '0')
		}
	}
	i, ok := nMap[s]
	if ok {
		return true, true, i
	}

	// check if string is substring of a number
	if len(s) > maxLen {
		return false, false, 0
	}
	for k := range nMap {
		if prefix && strings.HasPrefix(k, s)  || !prefix && strings.HasSuffix(k, s) {
			return false, true, 0
		}
	}
	return false, false, 0
}

func getLineSum(line string) (int, error) {
	fi := -1
	li := -1
	var fn int
	var ln int

	// for each starting index i, check all possible substrings
	for i := range line {
		for j := 1; j <= maxLen; j++ {
			if i+j > len(line) {
				break
			}
			s := line[i : i+j]
			isNum, maybe, n := isPossibleNumber(s, true)
            // found first digit
			if isNum {
                fi = i
                fn = n
                break
			}
			// break from inner loop and try next i
			if !maybe {
				break
			}
		}
        if fi != -1 {
            break
        }
	}

	for i := len(line); i >= 0; i-- {
		for j := 1; j <= maxLen; j++ {
            if i-j < 0 {
                break
            }
            s := line[i-j : i]
            isNum, maybe, n := isPossibleNumber(s, false)
            // found last digit
            if isNum {
                li = i
                ln = n
                break
            }
            // break from inner loop and try next i
            if !maybe {
                break
            }
        }
        if li != -1 {
            break
        }
	}

    // check if digits were found
	if fi == -1 || li == -1 {
		return 0, fmt.Errorf("No digits found")
	}

	// convert runes to int and return sum
	return fn * 10 + ln, nil
}

func calcTotal() (int, error) {
	// open file
	fpath := filepath.Join("data", file)
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file line by line, calculate sum
	total := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s, err := getLineSum(line)
		if err != nil {
			return 0, err
		}
		total += s
	}

	return total, nil
}

func main() {
	total, err := calcTotal()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Total: %d", total)
}
