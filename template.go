package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func main() {
	puzzleInput := LoadPuzzleInput()

	waitCount := 2
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(waitCount)

	partOneChannel := make(chan int)
	partTwoChannel := make(chan int)

	go DoPartOne(partOneChannel, waitGroup)
	go DoPartTwo(partTwoChannel, waitGroup)

	partOneResult := <- partOneChannel
	partTwoResult := <- partTwoChannel

	waitGroup.Wait()

	fmt.Printf("input: %v, part one: %d, part two: %d\n", puzzleInput, partOneResult, partTwoResult)
}

func DoPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 2
	waitGroup.Done()
}

func DoPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 1
	waitGroup.Done()
}

func LoadPuzzleInput() string {
	filename := "puzzle-input.dat"
	fd, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", filename, err))
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	return ""
}