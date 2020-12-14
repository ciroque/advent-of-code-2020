package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func (w *Window) isSumOfCandidates() (int64, bool) {
	candidateMap := w.buildMap()
	for _, number := range w.candidates {
		complement := w.target - number
		if _, exists := candidateMap[complement]; exists {
			return w.target, true
		}
	}

	return w.target, false
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	waitCount := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	partOneChannel := make(chan int64)
	partTwoChannel := make(chan int)

	go doExamples(&waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	log.Info().Int64("part-one", partOneResult).Int("part-two", partTwoResult).Msg("day ...")
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan int64, waitGroup *sync.WaitGroup) {
	const WindowSize = 25
	puzzleInput := loadPuzzleInput()
	numbers := buildNumberList(puzzleInput)
	windows := buildWindowList(WindowSize, numbers)

	for _, window := range windows {
		number, isSum := window.isSumOfCandidates()
		if !isSum {
			channel <- number
			break
		}
	}

	waitGroup.Done()
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

func doPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 2
	waitGroup.Done()
}

func loadPuzzleInput() []string {
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
