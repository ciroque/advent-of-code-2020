package main

import (
	"fmt"
	"github.com/ciroque/advent-of-code-2020/support"
	"sync"
)

func main() {
	puzzleInput := loadPuzzleInput()

	waitCount := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	partOneChannel := make(chan int)
	partTwoChannel := make(chan int)

	go doExamples(&waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	fmt.Printf("input: %v, part one: %d, part two: %d\n", puzzleInput, partOneResult, partTwoResult)
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 1
	waitGroup.Done()
}

func doPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 2
	waitGroup.Done()
}

func loadPuzzleInput() string {
	filename := "puzzle-input.dat"
	return support.ReadFile(filename)
}
