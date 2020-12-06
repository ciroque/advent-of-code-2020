package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"sync"
)

type Seat struct {
	row    int
	column int
}

func (s Seat) CalculateSeatId() int {
	return s.row*8 + s.column
}

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
	examples := []string{
		"BFFFBBFRRR",
		"FFFBBBFRRR",
		"BBFFBBFRLL",
	}

	for _, example := range examples {
		seat := ProcessBinarySpacePartitioning(example)
		fmt.Println(example, seat.row, seat.column, seat.CalculateSeatId())
	}

	waitGroup.Done()
}

func DoPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	highestId := 0
	puzzleInput := LoadPuzzleInput()

	for _, input := range puzzleInput {
		seat := ProcessBinarySpacePartitioning(input)
		seatId := seat.CalculateSeatId()
		if seatId > highestId {
			highestId = seatId
		}
	}

	channel <- highestId
	waitGroup.Done()
}

func DoPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	puzzleInput := LoadPuzzleInput()
	seatIds := BuildSeatIdsMap()

	for _, input := range puzzleInput {
		seat := ProcessBinarySpacePartitioning(input)
		seatId := seat.CalculateSeatId()
		seatIds[seatId] = true
	}

	var filledSeats []int
	for id, filled := range seatIds {
		if filled {
			filledSeats = append(filledSeats, id)
		}
	}

	sort.Ints(filledSeats)

	mySeatId := 0
	for index := 1; index < len(filledSeats); index++ {
		if filledSeats[index]-filledSeats[index-1] != 1 {
			mySeatId = filledSeats[index] - 1
			break
		}
	}

	channel <- mySeatId
	waitGroup.Done()
}

func ProcessBinarySpacePartitioning(seatSpecification string) Seat {
	front := func(slice []int) []int {
		return slice[0:(len(slice) / 2)]
	}

	back := func(slice []int) []int {
		return slice[(len(slice) / 2):]
	}

	rows := BuildRowsArray()
	columns := BuildColumnsArray()
	for _, char := range seatSpecification {
		if char == 'F' {
			rows = front(rows)
		} else if char == 'B' {
			rows = back(rows)
		} else if char == 'R' {
			columns = back(columns)
		} else if char == 'L' {
			columns = front(columns)
		}
	}

	return Seat{
		row:    rows[0],
		column: columns[0],
	}
}

func BuildSeatIdsMap() map[int]bool {
	rowCount := 128
	columnCount := 8

	seatIds := make(map[int]bool)

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
			seatIds[rowIndex*8+columnIndex] = false
		}
	}

	return seatIds
}

func BuildRowsArray() []int {
	rows := make([]int, 128)
	for index, _ := range rows {
		rows[index] = index
	}

	return rows
}

func BuildColumnsArray() []int {
	columns := make([]int, 8)
	for index, _ := range columns {
		columns[index] = index
	}

	return columns
}

func LoadPuzzleInput() []string {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	fd, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", filename, err))
	}

	content := []string{}

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		content = append(content, line)
	}

	err = fd.Close()
	if err != nil {
		fmt.Println(fmt.Errorf("error closing file: %s: %v", filename, err))
	}

	return content
}
