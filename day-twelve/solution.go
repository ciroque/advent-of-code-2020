package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"strconv"
	"sync"
	"time"
)

type Result struct {
	answer   int
	duration int64
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

	headingMap := map[int]rune{
		0:   'N',
		90:  'E',
		180: 'S',
		270: 'W',
	}

	north := 0
	east := 0
	heading := 90

	updateLocation := func(op rune, distance, n, e int) (n1, e1 int) {
		n1 = n
		e1 = e
		switch op {
		case 'N':
			n1 = n1 + distance
		case 'S':
			n1 = n1 - distance
		case 'E':
			e1 = e1 + distance
		case 'W':
			e1 = e1 - distance
		}

		return
	}

	for _, value := range puzzleInput {
		op := rune(value[0])
		wtf := value[1:]
		arg, _ := strconv.ParseInt(wtf, 10, 32)
		distance := int(arg)

		switch op {
		case 'L':
			heading = (heading - distance) % 360
			if heading > 360 || heading < 0 {
				heading = heading + 360
			}

		case 'R':
			heading = (heading + distance) % 360
			if heading > 360 {
				heading = heading - 360
			}

		case 'F':
			north, east = updateLocation(headingMap[heading], distance, north, east)

		default:
			north, east = updateLocation(op, distance, north, east)
		}

	}

	manhattanDistance := int(math.Abs(float64(north)) + math.Abs(float64(east)))

	channel <- Result{
		answer:   manhattanDistance,
		duration: time.Since(start).Microseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	channel <- Result{
		answer:   2,
		duration: time.Since(start).Nanoseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput() []string {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	return support.ReadFileIntoLines(filename)
}
