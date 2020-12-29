package main

import (
	"github.com/ciroque/advent-of-code-2020/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Result struct {
	answer   int
	duration int64
}

type PuzzleInput struct {
	ruleMap              map[string]Bounds
	ticket               string
	nearbyTickets        string
	parsedNearbyTickets  [][]int
	invalidTicketIndexes map[int]int
	errorRate            int
	validTickets         [][]int
}

type Bounds struct {
	upper Bound
	lower Bound
}

type Bound struct {
	upper int
	lower int
}

func (pi *PuzzleInput) findInvalidTickets() {
	const noLimit = -1
	const comma = ","
	values := strings.Split(strings.Replace(pi.nearbyTickets, "\n", ",", noLimit), comma)

	for ticketIndex, value := range values {
		i64val, _ := strconv.ParseInt(value, 10, 32)
		ival := int(i64val)

		inBounds := map[bool]int{
			true:  0,
			false: 0,
		}

		for _, bounds := range pi.ruleMap {
			inBounds[bounds.inRange(ival)]++
		}

		if inBounds[true] == 0 {
			pi.errorRate += ival
			pi.invalidTicketIndexes[ticketIndex] = ticketIndex
		}
	}

	return
}

func (pi *PuzzleInput) loadValidTickets() {

}

func (b *Bound) inRange(value int) bool {
	return value >= b.lower && value <= b.upper
}

func (bs *Bounds) inRange(value int) bool {
	return bs.lower.inRange(value) || bs.upper.inRange(value)
}

func newBounds(str string) (bounds Bounds) {
	parts := strings.Split(str, " or ")

	bounds = Bounds{
		upper: newBound(parts[1]),
		lower: newBound(parts[0]),
	}

	return
}

func newBound(str string) (bound Bound) {
	parts := strings.Split(str, "-")
	upper, _ := strconv.ParseInt(parts[1], 10, 32)
	lower, _ := strconv.ParseInt(parts[0], 10, 32)

	bound = Bound{
		upper: int(upper),
		lower: int(lower),
	}

	return
}

func newRuleMap(rules []string) (ruleMap map[string]Bounds) {
	ruleMap = map[string]Bounds{}

	for _, rule := range rules {
		parts := strings.Split(rule, ": ")

		ruleMap[parts[0]] = newBounds(parts[1])
	}

	return
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
		Msg("day sixteen")
}

func doExamples(waitGroup *sync.WaitGroup) {
	puzzleInput := loadPuzzleInput("example")
	puzzleInput.findInvalidTickets()
	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	puzzleInput := loadPuzzleInput("puzzle")
	puzzleInput.findInvalidTickets()

	channel <- Result{
		answer:   puzzleInput.errorRate,
		duration: time.Since(start).Milliseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	puzzleInput := loadPuzzleInput("puzzle")
	puzzleInput.findInvalidTickets()

	channel <- Result{
		answer:   len(puzzleInput.invalidTicketIndexes),
		duration: time.Since(start).Milliseconds(),
	}
	waitGroup.Done()
}

func loadPuzzleInput(prefix string) (puzzleInput PuzzleInput) {
	filename := prefix + "-input.dat"

	const comma = ","
	const newline = "\n"
	const noLimit = -1

	content := support.ReadFile(filename)
	sections := strings.Split(content, newline+newline) // [3]string
	ticket := strings.Split(sections[1], newline)

	justNearbyTickets := strings.Replace(sections[2], "nearby tickets:\n", "", noLimit)
	values := strings.Split(strings.Replace(justNearbyTickets, "\n", ",", noLimit), comma)

	var parsedNearbyTickets [][]int
	for _, value := range values {
		var nearbyTicket []int
		i64val, _ := strconv.ParseInt(value, 10, 32)
		ival := int(i64val)
		nearbyTicket = append(nearbyTicket, ival)
	}

	return PuzzleInput{
		ruleMap:              newRuleMap(strings.Split(sections[0], newline)),
		ticket:               ticket[1],
		nearbyTickets:        justNearbyTickets,
		parsedNearbyTickets:  parsedNearbyTickets,
		invalidTicketIndexes: map[int]int{},
	}
}
