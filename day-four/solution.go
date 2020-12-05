package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Passport struct {
	Fields map[string]string
	Valid  bool
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
	validPassportCount := 0
	passports := LoadPuzzleInput()
	validFields := []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}

	passportIsValid := func(passport *Passport) bool {
		for _, field := range validFields {
			if _, found := passport.Fields[field]; !found {
				return false
			}
		}
		return true
	}

	for _, passport := range passports {
		if passportIsValid(passport) {
			validPassportCount = validPassportCount + 1
		}
	}

	channel <- validPassportCount
	waitGroup.Done()
}

func DoPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- 2
	waitGroup.Done()
}

func LoadPuzzleInput() []*Passport {
	//filename := "example-input.dat"
	filename := "puzzle-input.dat"
	fd, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", filename, err))
	}

	passport := &Passport{Fields: map[string]string{}}
	passports := make([]*Passport, 0)
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			passports = append(passports, passport)
			passport = &Passport{Fields: map[string]string{}}
		} else {
			raw := strings.Split(line, " ")
			for _, wtf := range raw {
				parts := strings.Split(wtf, ":")
				passport.Fields[parts[0]] = parts[1]
			}
		}
	}
	passports = append(passports, passport)

	err = fd.Close()
	if err != nil {
		fmt.Println(fmt.Errorf("error closing file: %s: %v", filename, err))
	}

	return passports
}
