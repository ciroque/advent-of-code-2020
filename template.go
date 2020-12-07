package main

import (
	"fmt"
	"github.com/ciroque/advent-of-code-2020/support"
	"sync"
)

func main() {
	puzzleInput := LoadPuzzleInput()

	waitCount := 3
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(waitCount)

	partOneChannel := make(chan int)
	partTwoChannel := make(chan int)

	go DoExamples(waitGroup)
	go DoPartOne(partOneChannel, waitGroup)
	go DoPartTwo(partTwoChannel, waitGroup)

	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	fmt.Printf("input: %v, part one: %d, part two: %d\n", puzzleInput, partOneResult, partTwoResult)
}

func DoExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func DoPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 1
	waitGroup.Done()
}

func DoPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 2
	waitGroup.Done()
}

func LoadPuzzleInput() string {
	filename := "puzzle-input.dat"
	return support.ReadFile(filename)
}
