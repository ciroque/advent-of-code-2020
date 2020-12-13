package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
)

//var all map[string]map[string]int

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	waitCount := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	exampleChannel := make(chan int)
	partOneChannel := make(chan int)
	partTwoChannel := make(chan int)

	go doExamples(exampleChannel, &waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	exampleResult := <-exampleChannel
	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	log.Info().Int("example", exampleResult).Int("part-one", partOneResult).Int("part-two", partTwoResult).Msg("day seven")
}

func doExamples(channel chan int, waitGroup *sync.WaitGroup) {
	puzzleInput := loadExamplePuzzleInput()
	channel <- len(puzzleInput)
	waitGroup.Done()
}

func doPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	myBag := "shiny gold"
	puzzleInput := loadPuzzleInput()
	graph := buildGraph(puzzleInput, myBag)
	count := len(graph) - 1 // Should not include the `myBag`

	channel <- count
	waitGroup.Done()
}

func doPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	myBag := "shiny gold"
	puzzleInput := loadPuzzleInput()
	graph := buildGraph(puzzleInput, myBag)

	count := 0

	for _, v := range graph {
		count = count + len(v)
	}

	channel <- count
	waitGroup.Done()
}

func buildGraph(puzzleInput []string, myBag string) map[string]map[string]int {
	graph := map[string]map[string]int{}
	findAllContainers(puzzleInput, map[string]int{myBag: 1}, graph)
	return graph
}

func extractSubject(line string) string {
	separator := " "
	return strings.Join(strings.Split(line, separator)[0:2], separator)
}

func findContainers(puzzleInput []string, target string, graph map[string]map[string]int) map[string]int {
	if _, found := graph[target]; !found {
		graph[target] = map[string]int{}
	}

	containers := map[string]int{}
	for _, line := range puzzleInput {
		if foundAt := strings.Index(line, target); foundAt > 0 {
			subject := extractSubject(line)
			graph[target][subject]++
			containers[subject]++
		}
	}

	return containers
}

func findAllContainers(puzzleInput []string, containers map[string]int, graph map[string]map[string]int) {
	for container := range containers {
		nextContainers := findContainers(puzzleInput, container, graph)
		if len(nextContainers) > 0 {
			findAllContainers(puzzleInput, nextContainers, graph)
		}
	}
}

func loadExamplePuzzleInput() []string {
	filename := "example-input.dat"
	return support.ReadFileIntoLines(filename)
}

func loadPuzzleInput() []string {
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
