package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"sync"
)

type Instruction struct {
	Op         string
	Arg        int
	OriginalOp string
}

func (i *Instruction) swapOp() {
	if i.Op == "jmp" {
		i.Op = "nop"
	} else if i.Op == "nop" {
		i.Op = "jmp"
	}
}

func (i *Instruction) restoreOp() {
	i.Op = i.OriginalOp
}

func BuildInstructionList(puzzleInput []string) []Instruction {
	parseLine := func(line string) Instruction {
		parts := strings.Split(line, " ")
		argument, _ := strconv.ParseInt(parts[1], 0, 0)
		return Instruction{
			Op:         parts[0],
			Arg:        int(argument),
			OriginalOp: parts[0],
		}
	}

	var instructions []Instruction
	for _, line := range puzzleInput {
		instructions = append(instructions, parseLine(line))
	}
	return instructions
}

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
	instructions := BuildInstructionList(puzzleInput)
	_, accumulator, cycle := executeFrom(instructions)
	if cycle == true {
		log.Info().Msg("Cycle detected")
	}
	channel <- accumulator
	waitGroup.Done()
}

func doPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	puzzleInput := loadPuzzleInput()
	instructions := BuildInstructionList(puzzleInput)

	speculate := func(ops []int) (bool, int) {
		acc := 0
		foundProgramEnd := false
		for _, ip := range ops {
			instructions[ip].swapOp()
			_, acc, foundProgramEnd = executeFrom(instructions)
			if foundProgramEnd {
				return true, acc
			}
			instructions[ip].restoreOp()
		}

		return false, 0
	}

	foundProgramEnd := false
	var accumulator int
	nops := findOps(instructions, "nop")

	foundProgramEnd, accumulator = speculate(nops)
	if !foundProgramEnd {
		jmps := findOps(instructions, "jmp")
		foundProgramEnd, accumulator = speculate(jmps)
	}

	channel <- accumulator
	waitGroup.Done()
}

// returns int, int, bool
// The instruction pointer
// The accumulator
// true if a instruction length + 1 found, false otherwise
func executeFrom(instructions []Instruction) (int, int, bool) {
	instructionPointer := 0
	instructionCounts := map[int]int{}
	for index := range instructions {
		instructionCounts[index] = 0
	}

	accumulator := 0
	for true {
		instruction := instructions[instructionPointer]

		if instruction.Op == "acc" {
			accumulator = accumulator + instruction.Arg
			instructionPointer = instructionPointer + 1
		} else if instruction.Op == "jmp" {
			instructionPointer = instructionPointer + instruction.Arg
		} else {
			instructionPointer = instructionPointer + 1
		}

		instructionCounts[instructionPointer]++
		if instructionCounts[instructionPointer] > 1 {
			return instructionPointer, accumulator, false
		} else if instructionPointer == len(instructions) {
			return instructionPointer, accumulator, true
		}
	}

	return instructionPointer, accumulator, false
}

func findOps(instructions []Instruction, op string) []int {
	var ops []int
	for index, instruction := range instructions {
		if strings.Index(instruction.Op, op) == 0 {
			ops = append(ops, index)
		}
	}
	return ops
}

func loadPuzzleInput() []string {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
