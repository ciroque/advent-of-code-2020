package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

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

	log.Info().Int("part-one", partOneResult).Int("part-two", partTwoResult).Msg("day seven")
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	myGoldBag := "shiny gold"
	puzzleInput := loadPuzzleInput()

	bags := map[string]int{}
	for _, line := range puzzleInput {
		if len(line) > 0 {
			bags[extractSubject(line)]++
		}
	}

	outermostContainers := findContainers(puzzleInput, myGoldBag)

	topLevelBags := map[string]int{}
	count := findAllContainers(puzzleInput, outermostContainers, topLevelBags, bags)

	channel <- count
	waitGroup.Done()
}

func extractSubject(line string) string {
	separator := " "
	return strings.Join(strings.Split(line, separator)[0:2], separator)
}

func findContainers(puzzleInput []string, target string) map[string]int {
	containers := map[string]int{}
	for _, line := range puzzleInput {
		if foundAt := strings.Index(line, target); foundAt > 0 {
			containers[extractSubject(line)]++
		}
	}

	return containers
}

func findAllContainers(puzzleInput []string, containers map[string]int, topLevelBags map[string]int, bags map[string]int) int {
	count := 0

	for container, _ := range containers {
		nextContainers := findContainers(puzzleInput, container)
		if len(nextContainers) > 0 {
			count = count + findAllContainers(puzzleInput, nextContainers, topLevelBags, bags)
		} else {
			topLevelBags[container]++
			count = count + 1
		}
	}

	return count
}

func doPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 2
	waitGroup.Done()
}

func loadPuzzleInput() []string {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
