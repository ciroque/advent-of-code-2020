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
	answer   int64
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
		Int64("part-one-answer", partOneResult.answer).
		Int64("part-one-duration", partOneResult.duration).
		Int64("part-two-answer", partTwoResult.answer).
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
	sort.Slice(numbers, func(l, r int) bool { return numbers[l] < numbers[r] })

	deltas := map[int64]int{
		1: 1,
		3: 0,
	}
	for index := 1; index < len(numbers); index++ {
		delta := numbers[index] - numbers[index-1]
		deltas[delta]++
	}

	product := int64(deltas[1] * (deltas[3] + 1))

	channel <- Result{
		answer:   product,
		duration: time.Since(start).Microseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	areDifferentByOne := func(current, previous, next int64) bool {
		return current-previous == 1 && next-current == 1
	}

	puzzleInput := loadPuzzleInput()
	numbers := buildNumberList(puzzleInput)
	sort.Slice(numbers, func(l, r int) bool { return numbers[l] < numbers[r] })

	removables := make([]bool, len(numbers))
	removables[0] = areDifferentByOne(numbers[0], 0, numbers[1])
	for current, previous, next := 1, 0, 2; next < len(numbers); current, previous, next = current+1, previous+1, next+1 {
		removables[current] = areDifferentByOne(numbers[current], numbers[previous], numbers[next])
	}

	accumulator := int64(1)
	span := 0
	for _, value := range removables {
		if value {
			span = span + 1
		} else {
			switch span {
			case 1:
				accumulator *= 2
			case 2:
				accumulator *= 4
			case 3:
				accumulator *= 7
			default:
				log.Info().Int("span", span).Msg("ERROR: unrecognized span")
			}

			span = 0
		}
	}

	channel <- Result{
		answer:   accumulator,
		duration: time.Since(start).Microseconds(),
	}
	waitGroup.Done()
}

func buildNumberList(input []string) []int64 {
	var numbers []int64
	for _, line := range input {
		if number, err := strconv.ParseInt(line, 10, 64); err == nil {
			numbers = append(numbers, number)
		}
	}
	return numbers
}

func loadPuzzleInput() []string {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
