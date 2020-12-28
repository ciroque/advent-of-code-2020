package main

import (
	"fmt"
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Result struct {
	answer   int
	duration int64
}

type Location int

type State rune

const (
	Floor    State = '.'
	Occupied State = '#'
	Vacant   State = 'L'
)

const (
	TopLeft Location = iota
	BottomLeft
	TopRight
	BottomRight
	TopRow
	BottomRow
	LeftSide
	RightSide
	Body
)

type SeatingArea struct {
	rowCount         int
	rowWidth         int
	primarySeatMap   []State
	secondarySeatMap []State

	currentSeatMap  *[]State
	previousSeatMap *[]State
	changeCount     int
}

func (s *SeatingArea) areSeatingMapsSame() bool {
	for index := range s.primarySeatMap {
		if s.primarySeatMap[index] != s.secondarySeatMap[index] {
			return false
		}
	}

	return true
}

func (s *SeatingArea) chooseSeat(index int) (changed bool) {
	changed = false

	if s.isFloor(index) {
		changed = false
	} else if s.isPreviouslyOccupied(index) {
		if s.findAdjacentCount(index) >= 4 {
			(*s.currentSeatMap)[index] = Vacant
			changed = true
		} else {
			(*s.currentSeatMap)[index] = Occupied
			changed = false
		}
	} else { // seat is empty
		if s.findAdjacentCount(index) == 0 {
			(*s.currentSeatMap)[index] = Occupied
			changed = true
		} else {
			(*s.currentSeatMap)[index] = Vacant
			changed = false
		}
	}

	return
}

func (s *SeatingArea) chooseSeats() (changeCount int) {
	changeCount = 0

	for index := range *s.currentSeatMap {
		if changed := s.chooseSeat(index); changed {
			changeCount = changeCount + 1
		}
	}

	return
}

func (s *SeatingArea) countOccupiedSeats() (count int) {
	count = 0

	for index := range *s.currentSeatMap {
		if s.isOccupied(index) {
			count = count + 1
		}
	}

	return
}

func (s *SeatingArea) determineLocation(index int) (location Location) {
	location = Body

	row := index / s.rowWidth
	column := index % s.rowWidth

	if row == 0 && column == 0 {
		location = TopLeft

	} else if row == s.rowCount-1 && column == 0 {
		location = BottomLeft

	} else if row == 0 && column == s.rowWidth-1 {
		location = TopRight

	} else if row == s.rowCount-1 && column == s.rowWidth-1 {
		location = BottomRight

	} else if row == 0 {
		location = TopRow

	} else if row == s.rowCount-1 {
		location = BottomRow

	} else if column == 0 {
		location = LeftSide

	} else if column == s.rowWidth-1 {
		location = RightSide

	} else {
		location = Body
	}

	return
}

func (s *SeatingArea) findAdjacentCount(index int) (adjacentCount int) {
	adjacentCount = 0

	adjacentIndexes := s.findAdjacentIndexes(index)
	for _, adjacentIndex := range adjacentIndexes {
		if s.isPreviouslyOccupied(adjacentIndex) { // Needs to work on s.previousSeatMap
			adjacentCount = adjacentCount + 1
		}
	}

	return
}

func (s *SeatingArea) findAdjacentIndexes(index int) (adjacentIndexes []int) {
	nwi := index - s.rowWidth - 1
	ni := index - s.rowWidth
	nei := index - s.rowWidth + 1

	wi := index - 1
	ei := index + 1

	swi := index + s.rowWidth - 1
	si := index + s.rowWidth
	sei := index + s.rowWidth + 1

	location := s.determineLocation(index)
	switch location {
	case TopLeft:
		adjacentIndexes = []int{
			ei, si, sei,
		}
	case BottomLeft:
		adjacentIndexes = []int{
			ni, nei, ei,
		}
	case TopRight:
		adjacentIndexes = []int{
			wi, swi, si,
		}
	case BottomRight:
		adjacentIndexes = []int{
			nwi, ni, wi,
		}
	case TopRow:
		adjacentIndexes = []int{
			wi, ei, swi, si, sei,
		}
	case BottomRow:
		adjacentIndexes = []int{
			nwi, ni, nei, wi, ei,
		}
	case LeftSide:
		adjacentIndexes = []int{
			ni, nei, ei, si, sei,
		}
	case RightSide:
		adjacentIndexes = []int{
			nwi, ni, wi, swi, si,
		}
	case Body:
		adjacentIndexes = []int{
			nwi, ni, nei, wi, ei, swi, si, sei,
		}
	}

	return
}

func (s *SeatingArea) isFloor(index int) bool {
	return (*s.currentSeatMap)[index] == Floor
}

func (s *SeatingArea) isOccupied(index int) bool {
	return (*s.currentSeatMap)[index] == Occupied
}

func (s *SeatingArea) isPreviouslyOccupied(index int) bool {
	return (*s.previousSeatMap)[index] == Occupied
}

func (s *SeatingArea) load(input []string) {
	s.rowCount = len(input)
	s.rowWidth = len(input[0])
	s.primarySeatMap = make([]State, s.rowCount*s.rowWidth)
	s.secondarySeatMap = make([]State, s.rowCount*s.rowWidth)

	s.currentSeatMap = &s.primarySeatMap
	s.previousSeatMap = &s.secondarySeatMap

	for lineIndex, line := range input {
		for runeIndex, char := range line {
			index := (lineIndex * s.rowWidth) + (runeIndex % s.rowWidth)
			(*s.currentSeatMap)[index] = State(char)
			(*s.previousSeatMap)[index] = State(char)
		}
	}
}

func (s *SeatingArea) printCurrentMap() {
	fmt.Println()
	fmt.Println("current map")
	for index, value := range *s.currentSeatMap {
		fmt.Print(string(value))
		if (index+1)%s.rowWidth == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func (s *SeatingArea) printPreviousMap() {
	fmt.Println()
	fmt.Println("previous map")
	for index, value := range *s.previousSeatMap {
		fmt.Print(string(value))
		if (index+1)%s.rowWidth == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func (s *SeatingArea) swapSeatMaps() {
	temp := s.currentSeatMap
	s.currentSeatMap = s.previousSeatMap
	s.previousSeatMap = temp
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	waitCount := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(waitCount)

	partOneChannel := make(chan Result)
	partTwoChannel := make(chan Result)

	go doExamples(&waitGroup)
	go doPartOne(partOneChannel, &waitGroup)
	go doPartTwo(partTwoChannel, &waitGroup)

	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	log.
		Info().
		Int("part-one-answer", partOneResult.answer).
		Int64("part-one-duration", partOneResult.duration).
		Int("part-two-answer", partTwoResult.answer).
		Int64("part-two-duration", partTwoResult.duration).
		Msg("day ...")
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	puzzleInput := loadPuzzleInput()

	seatingArea := SeatingArea{}
	seatingArea.load(puzzleInput)

	changeCount := seatingArea.chooseSeats()
	//seatingArea.printCurrentMap()
	for changeCount > 0 {
		seatingArea.swapSeatMaps()
		changeCount = seatingArea.chooseSeats()
		//seatingArea.printCurrentMap()
	}

	channel <- Result{
		answer:   seatingArea.countOccupiedSeats(),
		duration: time.Since(start).Milliseconds(),
	}

	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   1,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput() []string {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
