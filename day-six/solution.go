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
	puzzleInput := LoadPuzzleInput()
	parsedPuzzleInput := strings.Split(puzzleInput, "\n")

	processResponse := func(response string, answers map[int32]int) {
		for _, question := range response {
			if _, found := answers[question]; found {
				answers[question] = answers[question] + 1
			} else {
				answers[question] = 1
			}
		}
	}

	checkResponses := func(count int, answers map[int32]int) int {
		sum := 0
		for _, v := range answers {
			if v == count {
				sum = sum + 1
			}
		}

		return sum
	}

	sum := 0
	count := 0
	answers := map[int32]int{}
	for _, response := range parsedPuzzleInput {
		if len(response) == 0 {
			sum = sum + checkResponses(count, answers)
			answers = map[int32]int{}
			count = 0
		} else {
			count = count + 1
			processResponse(response, answers)
		}
	}

	channel <- sum
	waitGroup.Done()
}

func LoadPuzzleInput() string {
	filename := "puzzle-input.dat"
	return support.ReadFile(filename)
}
