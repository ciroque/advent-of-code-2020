package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
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

	log.Info().Int("part-one", partOneResult).Int("part-two", partTwoResult).Msg("day 8")
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	puzzleInput := loadPuzzleInput()

	instructions := map[int]int{}
	accumulator := 0
	for index := range puzzleInput {
		instructions[index] = 0
	}

	instructionPointer := 0
	for true {
		line := puzzleInput[instructionPointer]
		op, arg := parseLine(line)

		if op == "acc" {
			accumulator = accumulator + arg
			instructionPointer = instructionPointer + 1
		} else if op == "jmp" {
			instructionPointer = instructionPointer + arg
		} else {
			instructionPointer = instructionPointer + 1
		}

		instructions[instructionPointer]++
		if instructions[instructionPointer] > 1 {
			break
		}
	}

	channel <- accumulator
	waitGroup.Done()
}

func doPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 2
	waitGroup.Done()
}

func parseLine(line string) (string, int) {
	parts := strings.Split(line, " ")
	argument, _ := strconv.ParseInt(parts[1], 0, 0)
	return parts[0], int(argument)
}

func loadPuzzleInput() []string {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
