package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

	RunTestValues()

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

	for _, passport := range passports {
		if HasRequiredFields(passport) {
			validPassportCount = validPassportCount + 1
		}
	}

	channel <- validPassportCount
	waitGroup.Done()
}

func DoPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	validPassportCount := 0
	passports := LoadPuzzleInput()

	for _, passport := range passports {
		hasFields := HasRequiredFields(passport)
		hasValues := HasValidValues(passport)
		if hasFields && hasValues {
			validPassportCount = validPassportCount + 1
		}
	}

	channel <- validPassportCount
	waitGroup.Done()
}

func HasRequiredFields(passport *Passport) bool {
	validFields := []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}

	for _, field := range validFields {
		if _, found := passport.Fields[field]; !found {
			return false
		}
	}
	return true
}

func HasValidValues(passport *Passport) bool {
	validators := BuildValidationMap()

	for key, value := range passport.Fields {
		if !validators[key](value) {
			fmt.Println("failed", key, value)
			return false
		}
	}

	return true
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

func BuildValidationMap() map[string]func(string) bool {
	validators := map[string]func(string) bool{}

	validators["byr"] = func(s string) bool {
		year, _ := strconv.Atoi(s)
		return year >= 1920 && year <= 2002
	}

	validators["iyr"] = func(s string) bool {
		year, _ := strconv.Atoi(s)
		return year >= 2010 && year <= 2020
	}

	validators["eyr"] = func(s string) bool {
		year, _ := strconv.Atoi(s)
		return year >= 2020 && year <= 2030
	}

	validators["hgt"] = func(s string) bool {
		if strings.Contains(s, "cm") {
			height, _ := strconv.Atoi(strings.TrimSuffix(s, "cm"))
			return height >= 150 && height <= 193
		} else if strings.Contains(s, "in") {
			height, _ := strconv.Atoi(strings.TrimSuffix(s, "in"))
			return height >= 59 && height <= 76
		} else {
			return false
		}
	}

	validators["hcl"] = func(s string) bool {
		regex := regexp.MustCompile("#[0-9a-f]{6}")
		return regex.Match([]byte(s))
	}

	validators["ecl"] = func(s string) bool {
		validValues := map[string]string{"amb": "amb", "blu": "blu", "brn": "brn", "gry": "gry", "grn": "grn", "hzl": "hzl", "oth": "oth"}
		_, found := validValues[s]
		return found
	}

	validators["pid"] = func(s string) bool {
		regex := regexp.MustCompile("[0-9]")
		return len(s) == 9 && regex.Match([]byte(s))
	}

	validators["cid"] = func(s string) bool {
		return true
	}

	return validators
}

func RunTestValues() bool {
	validators := BuildValidationMap()
	valid := false

	valid = validators["byr"]("2002") //valid
	valid = validators["byr"]("2003") // invalid

	valid = validators["hgt"]("60in")  // valid
	valid = validators["hgt"]("190cm") // valid
	valid = validators["hgt"]("190in") // invalid
	valid = validators["hgt"]("190")   // invalid

	valid = validators["hcl"]("#123abc") // valid
	valid = validators["hcl"]("#123abz") // invalid
	valid = validators["hcl"]("123abc")  // invalid

	valid = validators["ecl"]("brn") // valid
	valid = validators["ecl"]("wat") // invalid

	valid = validators["pid"]("000000001")  // valid
	valid = validators["pid"]("0123456789") // invalid

	return valid
}
