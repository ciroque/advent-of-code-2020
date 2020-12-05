package main

import (
  "bufio"
  "fmt"
  "os"
  "sort"
  "strconv"
  "sync"
)

func main() {
  goal := 2020
  data := ReadInputFile()

  waitCount := 2
  waitGroup := &sync.WaitGroup{}
  waitGroup.Add(waitCount)

  partOneChannel := make(chan int)
  partTwoChannel := make(chan int)

  go DoPartOne(data, goal, partOneChannel, waitGroup)
  go DoPartTwo(data, goal, partTwoChannel, waitGroup)

  partOneResult := <- partOneChannel
  partTwoResult := <- partTwoChannel

  waitGroup.Wait()

  fmt.Printf("part one: %d, part two: %d\n", partOneResult, partTwoResult)
}

func ReadInputFile() []int {
  filename := "puzzle-input.dat"
  fd, err := os.Open(filename)
  if err != nil {
    panic(fmt.Sprintf("open %s: %v", filename, err))
  }
  defer fd.Close()
  var numbers []int
  scanner := bufio.NewScanner(fd)
  for scanner.Scan() {
    lineStr := scanner.Text()
    number, _ := strconv.Atoi(lineStr)
    numbers = append(numbers, number)
  }

  return numbers
}

func DoPartOne(numbers []int, goal int, partOneChannel chan int, waitGroup *sync.WaitGroup) {
  numberMap := make(map[int]int)

  for _, number := range numbers {
    numberMap[number] = number
  }

  first, second, found := FindPair(numberMap, goal)
  if !found {
    message := fmt.Errorf("there was no pair in the list that meet the goal of %d", goal)
    fmt.Println(message)
  } else {
    fmt.Printf("The pair goal of %d was met by %d and %d whose product is %d\n", goal, first, second, first * second)
  }

  partOneChannel <- first * second

  waitGroup.Done()
}

func FindPair(numbers map[int]int, goal int) (first int, second int, found bool) {
  for number, _ := range numbers {
    complement := goal - number
    _, exists := numbers[complement]
    if exists {
      return number, complement, true
    }
  }

  return 0, 0, false
}

func DoPartTwo(numbers []int, goal int, partTwoChannel chan int, waitGroup *sync.WaitGroup) {
  sortedNumbers := make([]int, len(numbers))
  copy(sortedNumbers, numbers)
  sort.Ints(sortedNumbers)

  first, second, third, found := FindTriple(sortedNumbers, goal)
  if !found {
   message := fmt.Errorf("there was no triple in the list that meet the goal of %d", goal)
   fmt.Println(message)
  }
  fmt.Printf("The triple goal of %d was met by %d, %d, and %d whose product is %d\n", goal, first, second, third, first * second * third)

  partTwoChannel <- first * second * third
  waitGroup.Done()
}

func FindTriple(numbers []int, goal int) (int, int, int, bool) {
  for index, number := range numbers {
    lowIndex := index + 1
    highIndex := len(numbers) - 1

    for lowIndex < highIndex {
      sum := number + numbers[lowIndex] + numbers[highIndex]
      if sum == goal  {
        return number, numbers[lowIndex], numbers[highIndex], true
      } else if sum < goal {
        lowIndex = lowIndex + 1
      } else {
        highIndex = highIndex - 1
      }

      //fmt.Println(index, number, lowIndex, highIndex)
    }
  }
  return 0, 0, 0, false
}
