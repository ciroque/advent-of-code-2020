package main

import (
	"fmt"
	"github.com/ciroque/advent-of-code-2020/support"
	"regexp"
	"strings"
	"sync"
)

func main() {

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

	fmt.Printf("part one: %d, part two: %d\n", partOneResult, partTwoResult)
}

func DoExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func DoPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	puzzleInput := LoadPuzzleInput()
	groupPartitioningRegex := regexp.MustCompile("\n\n")
	groupCoalescingRegex := regexp.MustCompile("\n")

	partitionedGroups := groupPartitioningRegex.ReplaceAllString(puzzleInput, "\t")
	coalescedGroups := groupCoalescingRegex.ReplaceAllString(partitionedGroups, "")
	groupResponses := strings.Split(coalescedGroups, "\t")

	sum := 0
	for _, line := range groupResponses {
		mappedGroupResponses := map[int32]int32{}
		for _, char := range line {
			mappedGroupResponses[char] = char
		}
		sum = sum + len(mappedGroupResponses)
	}

	channel <- sum
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
