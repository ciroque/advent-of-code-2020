package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sort"
	"strconv"
	"sync"
)

type Window struct {
	target     int64
	candidates []int64
}

func (w *Window) buildMap() map[int64]int64 {
	candidateMap := map[int64]int64{}
	for _, candidate := range w.candidates {
		candidateMap[candidate] = candidate
	}
	return candidateMap
}

func (w *Window) hasSumOfPairs() (int64, bool) {
	candidateMap := w.buildMap()
	for _, number := range w.candidates {
		complement := w.target - number
		if _, exists := candidateMap[complement]; exists {
			return w.target, true
		}
	}

	return w.target, false
}

func (w *Window) hasSumOfContiguousNumbers() (int64, bool) {
	length := len(w.candidates)
	for startIndex := 0; startIndex < length; startIndex++ {
		accumulator := w.candidates[startIndex]
		for seekIndex := startIndex + 1; seekIndex < length; seekIndex++ {
			accumulator = accumulator + w.candidates[seekIndex]
			if accumulator == 31161678 {
				contiguousNumbers := w.candidates[startIndex:seekIndex]
				sort.Slice(contiguousNumbers, func(l, r int) bool { return contiguousNumbers[l] < contiguousNumbers[r] })
				sum := contiguousNumbers[0] + contiguousNumbers[len(contiguousNumbers)-1]

				return sum, true
			}
		}
	}

	return -1, false
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	waitCount := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	partOneChannel := make(chan int64)
	partTwoChannel := make(chan int64)

	go doExamples(&waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	log.Info().Int64("part-one", partOneResult).Int64("part-two", partTwoResult).Msg("day ...")
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan int64, waitGroup *sync.WaitGroup) {
	const WindowSize = 25
	puzzleInput := loadPuzzleInput()
	numbers := buildNumberList(puzzleInput)
	windows := buildWindowList(WindowSize, numbers)

	if number, found := findNonSumValue(windows); found {
		channel <- number
	} else {
		channel <- -1
	}

	waitGroup.Done()
}

func doPartTwo(channel chan int64, waitGroup *sync.WaitGroup) {
	const WindowSize = 25
	puzzleInput := loadPuzzleInput()
	numbers := buildNumberList(puzzleInput)
	windows := buildWindowList(WindowSize, numbers)

	if number, found := findContiguousSum(windows); found {
		channel <- number
	} else {
		channel <- -1
	}

	waitGroup.Done()
}

func findContiguousSum(windows []Window) (int64, bool) {
	for _, window := range windows {
		if sum, found := window.hasSumOfContiguousNumbers(); found {
			return sum, true
		}
	}
	return -1, false
}

func findNonSumValue(windows []Window) (int64, bool) {
	for _, window := range windows {
		number, isSum := window.hasSumOfPairs()
		if !isSum {
			return number, true
		}
	}
	return -1, false
}

func buildWindowList(WindowSize int, numbers []int64) []Window {
	var windows []Window
	numbersLen := len(numbers)
	for currentIndex := WindowSize; currentIndex < numbersLen; currentIndex++ {
		target := numbers[currentIndex]
		candidates := numbers[currentIndex-WindowSize : currentIndex]
		window := Window{
			target:     target,
			candidates: candidates,
		}
		windows = append(windows, window)
	}
	return windows
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
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
