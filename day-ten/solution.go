package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Result struct {
	answer   int
	duration int64
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	waitCount := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	partOneChannel := make(chan Result)
	partTwoChannel := make(chan Result)

	go doExamples(&waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	log.
		Info().
		Int("part-one-answer", partOneResult.answer).
		Int64("part-one-duration", partOneResult.duration).
		Int("part-two-answer", partTwoResult.answer).
		Int64("part-two-duration", partTwoResult.duration).
		Msg("day 10")
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()
	puzzleInput := loadPuzzleInput()
	numbers := buildNumberList(puzzleInput)
	sort.Ints(numbers)

	deltas := map[int]int{
		1: 1,
		3: 0,
	}
	for index := 1; index < len(numbers); index++ {
		delta := numbers[index] - numbers[index-1]
		deltas[delta]++
	}

	product := deltas[1] * (deltas[3] + 1)

	channel <- Result{
		answer:   product,
		duration: time.Since(start).Microseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   2,
		duration: time.Since(start).Microseconds(),
	}
	waitGroup.Done()
}

func buildNumberList(input []string) []int {
	var numbers []int
	for _, line := range input {
		if number, err := strconv.ParseInt(line, 10, 32); err == nil {
			numbers = append(numbers, int(number))
		}
	}
	return numbers
}

func loadPuzzleInput() []string {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
