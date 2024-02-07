package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// var file = "example.txt"
var file = "input.txt"
var bag = map[string]int{
    "red": 12, "green": 13, "blue": 14,
}

type Game struct {
    GameNum int
    Rounds []map[string]int
}

func newGame(n int) Game {
    return Game{
        GameNum: n, 
        Rounds: make([]map[string]int, 1),
    }
}

// addRound adds a round to the game
// s: "1 red, 2 green, 3 blue"
func (g *Game) addRound(s string) {
    round := map[string]int{}
    for _, entry := range strings.Split(s, ", ") {
        arr := strings.Split(entry, " ")
        color := arr[1]
        num, err := strconv.Atoi(arr[0])
        if err != nil {
            log.Fatal(err)
        }
        round[color] = num
    }
    g.Rounds = append(g.Rounds, round)
}

func (g *Game) getRequiredBalls() map[string]int {
    minBalls := map[string] int {
        "red": 0, "green": 0, "blue": 0,
    }
    for _, r := range g.Rounds {
        for color, num := range r {
            if num > minBalls[color] {
                minBalls[color] = num
            }
        }
    }
    return minBalls
}

func (g *Game) isPossible() bool {
    reqBalls := g.getRequiredBalls()
    for color, num := range reqBalls {
        if num > bag[color] {
            return false
        }
    }
    return true
}

func (g *Game) getPower() int {
    reqBalls := g.getRequiredBalls()
    power := 1 // assuming at least 1 ball of each color
    for _, num := range reqBalls {
        power *= num
    }
    return power
}

func parseGameLine(s string) Game {
    // split into game number and rounds
    parts := strings.Split(s, ": ")

    // game number
    gameNumStr := strings.Split(parts[0], " ")[1]
    gameNum, err := strconv.Atoi(gameNumStr)
    if err != nil {
        log.Fatal(err)
    }

    // create game and populate rounds
    g := newGame(gameNum)
    for _, rs := range strings.Split(parts[1], "; ") {
        g.addRound(rs)
    }

    return g
}

func main() {
    fpath := filepath.Join("data", file)
    file, err := os.Open(fpath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    games := make([]Game, 0)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        g := parseGameLine(line)
        games = append(games, g)
    }

    validGameNums := make([]int, 0)
    for _, g := range games {
        if g.isPossible() {
            validGameNums = append(validGameNums, g.GameNum)
        }
    }

    fmt.Printf("Valid game numbers: %v\n", validGameNums)
    fmt.Printf("Number of valid games: %d\n", len(validGameNums))

    // part 1: sum of valid game numbers
    sum := 0
    for _, n := range validGameNums {
        sum += n
    }
    fmt.Printf("P1. Sum of valid game numbers: %d\n", sum)

    // part 2: sum of powers
    powerSum := 0
    for _, g := range games {
        powerSum += g.getPower()
    }
    fmt.Printf("P2. Sum of powers: %d\n", powerSum)
}
