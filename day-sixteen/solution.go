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
	valuesCount          int
	ruleMap              map[string]Bounds
	rawNearbyTickets     string
	errorRate            int
	ticket               []int
	nearbyTickets        [][]int
	invalidTicketIndexes map[int]int
	validTickets         [][]int
	valueIndexMap        map[string]int
	departureProduct     int
}

type Bounds struct {
	upper Bound
	lower Bound
}

type Bound struct {
	upper int
	lower int
}

func (pi *PuzzleInput) calculateDepartureProduct() {
	for ruleName, valueIndex := range pi.valueIndexMap {
		if strings.HasPrefix(ruleName, "departure") {
			pi.departureProduct = pi.departureProduct * pi.ticket[valueIndex]
		}
	}
}

func (pi *PuzzleInput) calculateErrorRate() {
	const noLimit = -1
	const comma = ","
	values := strings.Split(strings.Replace(pi.rawNearbyTickets, "\n", ",", noLimit), comma)

	for valueIndex, value := range values {
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
			pi.invalidTicketIndexes[valueIndex/pi.valuesCount] = ival // total side-effect, but hell, we're here
		}
	}

	return
}

func (pi *PuzzleInput) cloneRuleMap() (ruleMap map[string]Bounds) {
	ruleMap = map[string]Bounds{}
	for ruleName, bounds := range pi.ruleMap {
		ruleMap[ruleName] = bounds
	}
	return
}

func (pi *PuzzleInput) loadValidTickets() {
	for index, _ := range pi.nearbyTickets {
		if _, found := pi.invalidTicketIndexes[index]; !found {
			pi.validTickets = append(pi.validTickets, pi.nearbyTickets[index])
		}
	}
}

func (pi *PuzzleInput) findValueFieldsIndexes() {
	ruleMap := pi.cloneRuleMap()
	resolvedValueIndexes := map[int]bool{}
	for len(pi.valueIndexMap) < len(pi.ruleMap) {

		for valueIndex := range pi.ticket {
			if resolvedValueIndexes[valueIndex] {
				continue
			}

			var validRuleNames []string
			for ruleName, bounds := range ruleMap {
				satisfiesAllBounds := true
				for _, ticket := range pi.validTickets {
					value := ticket[valueIndex]
					passed := bounds.inRange(value)
					if !passed {
						satisfiesAllBounds = false
						break
					}
				}

				if satisfiesAllBounds {
					validRuleNames = append(validRuleNames, ruleName)
				}
			}

			if len(validRuleNames) == 1 {
				ruleName := validRuleNames[0]
				pi.valueIndexMap[ruleName] = valueIndex
				resolvedValueIndexes[valueIndex] = true
				delete(ruleMap, ruleName)
			}
		}
	}
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
	puzzleInput := loadPuzzleInput("part-two")
	puzzleInput.calculateErrorRate()
	puzzleInput.loadValidTickets()
	puzzleInput.findValueFieldsIndexes()
	puzzleInput.calculateDepartureProduct()

	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	puzzleInput := loadPuzzleInput("puzzle")
	puzzleInput.calculateErrorRate()

	channel <- Result{
		answer:   puzzleInput.errorRate,
		duration: time.Since(start).Milliseconds(),
	}
	waitGroup.Done()
}

func doPartTwo(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	puzzleInput := loadPuzzleInput("puzzle")
	puzzleInput.calculateErrorRate()
	puzzleInput.loadValidTickets()
	puzzleInput.findValueFieldsIndexes()
	puzzleInput.calculateDepartureProduct()

	channel <- Result{
		answer:   puzzleInput.departureProduct,
		duration: time.Since(start).Milliseconds(),
	}
	waitGroup.Done()
}

func stringToIntList(data string) (numbers []int) {
	numbers = stringListToIntList(strings.Split(data, ","))
	return
}

func stringListToIntList(strings []string) (numbers []int) {
	for _, entry := range strings {
		number64, _ := strconv.ParseInt(entry, 10, 32)
		numbers = append(numbers, int(number64))
	}

	return
}

func loadPuzzleInput(prefix string) (puzzleInput PuzzleInput) {
	filename := prefix + "-input.dat"

	const comma = ","
	const newline = "\n"
	const noLimit = -1

	content := support.ReadFile(filename)
	sections := strings.Split(content, newline+newline)
	ticketParts := strings.Split(sections[1], newline)
	rawNearbyTickets := strings.Replace(sections[2], "nearby tickets:\n", "", noLimit)

	ticket := stringToIntList(ticketParts[1])

	var nearbyTickets [][]int
	for _, ticket := range strings.Split(rawNearbyTickets, "\n") {
		nearbyTickets = append(nearbyTickets, stringToIntList(ticket))
	}

	return PuzzleInput{
		valuesCount:          len(ticket),
		ruleMap:              newRuleMap(strings.Split(sections[0], newline)),
		ticket:               ticket,
		rawNearbyTickets:     rawNearbyTickets,
		nearbyTickets:        nearbyTickets,
		invalidTicketIndexes: map[int]int{},
		validTickets:         [][]int{},
		valueIndexMap:        map[string]int{},
		departureProduct:     1,
	}
}
