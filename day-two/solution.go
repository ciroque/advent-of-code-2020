package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type PasswordData struct {
	MinOccurrences    int
	MaxOccurrences    int
	RequiredCharacter int32
	Password          string
	Valid             bool
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
	example1 :=
		&PasswordData{
			MinOccurrences:    1,
			MaxOccurrences:    3,
			RequiredCharacter: 'a',
			Password:          "abcde",
			Valid:             false,
		}
	example2 := &PasswordData{
		MinOccurrences:    1,
		MaxOccurrences:    3,
		RequiredCharacter: 'b',
		Password:          "cdefg",
		Valid:             false,
	}

	example3 := &PasswordData{
		MinOccurrences:    2,
		MaxOccurrences:    9,
		RequiredCharacter: 'c',
		Password:          "ccccccccc",
		Valid:             false,
	}

	ValidateTobogganRentalPassword(example1)
	ValidateTobogganRentalPassword(example2)
	ValidateTobogganRentalPassword(example3)

	fmt.Fprintf(os.Stdout, "\n\texample1: { expected: true, actual: %t }\n\texample2: { expected: false, actual: %t }\n\texample3: { expected: false, actual: %t }\n\n", example1.Valid, example2.Valid, example3.Valid)

	waitGroup.Done()
}

func DoPartOne(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- ValidateAndCountPasswords(ValidateSledRentalPassword)
	waitGroup.Done()
}

func DoPartTwo(channel chan int, waitGroup *sync.WaitGroup) {
	channel <- ValidateAndCountPasswords(ValidateTobogganRentalPassword)
	waitGroup.Done()
}

func ValidateAndCountPasswords(validator func(data *PasswordData)) int {
	input := LoadPuzzleInput()
	for _, passwordData := range input {
		validator(passwordData)
	}

	count := 0
	for _, passwordData := range input {
		if passwordData.Valid {
			count = count + 1
		}
	}

	return count
}

func ValidateSledRentalPassword(passwordData *PasswordData) {
	passwordCharacterMap := make(map[int32]int)
	for _, character := range passwordData.Password {
		_, found := passwordCharacterMap[character]
		if found {
			passwordCharacterMap[character] = passwordCharacterMap[character] + 1
		} else {
			passwordCharacterMap[character] = 1
		}
	}

	passwordData.Valid = passwordCharacterMap[passwordData.RequiredCharacter] <= passwordData.MaxOccurrences && passwordCharacterMap[passwordData.RequiredCharacter] >= passwordData.MinOccurrences
}

func ValidateTobogganRentalPassword(passwordData *PasswordData) {
	requiredCharacter := uint8(passwordData.RequiredCharacter)
	firstPosition := passwordData.Password[passwordData.MinOccurrences-1] == requiredCharacter
	secondPosition := passwordData.Password[passwordData.MaxOccurrences-1] == requiredCharacter
	passwordData.Valid = secondPosition != firstPosition
}

func LoadPuzzleInput() []*PasswordData {
	passwordData := make([]*PasswordData, 0)
	filename := "puzzle-input.dat"
	fd, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", filename, err))
	}

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		passwordData = append(passwordData, ParsePuzzleInputLine(scanner.Text()))
	}

	err = fd.Close()
	if err != nil {
		fmt.Println(fmt.Errorf("error closing file: %s: %v", filename, err))
	}

	return passwordData
}

func ParsePuzzleInputLine(line string) *PasswordData {
	separator := " "
	occurrenceRangeIndex := 0
	characterIndex := 1
	passwordIndex := 2

	components := strings.Split(line, separator)

	occurrenceRange := strings.Split(components[occurrenceRangeIndex], "-")
	character := components[characterIndex][0]
	password := components[passwordIndex]

	minOccurrences, _ := strconv.Atoi(occurrenceRange[0])
	maxOccurrences, _ := strconv.Atoi(occurrenceRange[1])

	return &PasswordData{
		MinOccurrences:    minOccurrences,
		MaxOccurrences:    maxOccurrences,
		RequiredCharacter: int32(character),
		Password:          password,
	}
}
