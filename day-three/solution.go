package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type Grid struct {
	IsTree bool
}

func main() {
	waitCount := 2
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(waitCount)

	partOneChannel := make(chan int)
	partTwoChannel := make(chan int)

	go DoPartOne(partOneChannel, waitGroup)
	go DoPartTwo(partTwoChannel, waitGroup)

	partOneResult := <-partOneChannel
	partTwoResult := <-partTwoChannel

	waitGroup.Wait()

	fmt.Printf("part one: %d, part two: %d\n", partOneResult, partTwoResult)
}

func DoPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	const SLOPE_RIGHT = 3
	const SLOPE_DOWN = 1

	channel <- CountTreesOnSlope(SLOPE_RIGHT, SLOPE_DOWN)
	waitGroup.Done()
}

func DoPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	treeCountOne := CountTreesOnSlope(1, 1)
	treeCountTwo := CountTreesOnSlope(3, 1)
	treeCountThree := CountTreesOnSlope(5, 1)
	treeCountFour := CountTreesOnSlope(7, 1)
	treeCountFive := CountTreesOnSlope(1, 2)

	channel <- treeCountOne * treeCountTwo * treeCountThree * treeCountFour * treeCountFive
	waitGroup.Done()
}

func CountTreesOnSlope(slopeRight int, slopeDown int) int {
	width, height, grid := LoadPuzzleInput()

	index := 0
	row := 0
	column := 0
	treeCount := 0
	for row < height-1 {
		row = row + slopeDown
		column = column + slopeRight
		adjustedColumn := column % width
		index = (row * width) + adjustedColumn

		if grid[index].IsTree {
			treeCount = treeCount + 1
		}
	}

	return treeCount
}

func LoadPuzzleInput() (int, int, []Grid) {
	var grid []Grid
	var width int
	var height int
	filename := "puzzle-input.dat"
	fd, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", filename, err))
	}

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		width = len(line)
		height = height + 1
		for _, character := range line {
			gridPoint := Grid{
				IsTree: character == '#',
			}
			grid = append(grid, gridPoint)
		}
	}

	err = fd.Close()
	if err != nil {
		fmt.Println(fmt.Errorf("error closing file: %s: %v", filename, err))
	}

	return width, height, grid
}

func LoadExampleData() (int, int, []Grid) {
	line := "..##.......#...#...#...#....#..#...#.#...#.#.#...##..#...#.##......#.#.#....#.#........##.##...#...#...##....#.#..#...#.#"
	var grid []Grid

	for _, character := range line {
		gridPoint := Grid{
			IsTree: character == '#',
		}
		grid = append(grid, gridPoint)
	}

	return 11, 11, grid
}
