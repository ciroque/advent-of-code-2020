package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Result struct {
	answer   int64
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
		Int64("part-one-answer", partOneResult.answer).
		Int64("part-one-duration", partOneResult.duration).
		Int64("part-two-answer", partTwoResult.answer).
		Int64("part-two-duration", partTwoResult.duration).
		Msg("day ...")
}

func applyTransform(input int64, subject int64) int64 {
	return (input * subject) % 20201227
}

func calculateEncryptionKey(publicKey int64, loopSize int64) (encryptionKey int64) {
	encryptionKey = 1
	for index := int64(0); index < loopSize; index++ {
		encryptionKey = applyTransform(encryptionKey, publicKey)
	}

	return
}

func findLoopCount(target int64) int64 {
	const Subject = 7
	guess := int64(1)
	for index := int64(1); true; index++ {
		guess = applyTransform(guess, Subject)
		if guess == target {
			return index
		}
	}

	return -1
}

func doExamples(waitGroup *sync.WaitGroup) {

	waitGroup.Done()
}

func doPartOne(channel chan Result, waitGroup *sync.WaitGroup) {
	start := time.Now()

	cardPublicKey, doorPublicKey := loadPuzzleInput()
	cardLoopSize := findLoopCount(cardPublicKey)
	doorLoopSize := findLoopCount(doorPublicKey)

	cardEncryptionKey := calculateEncryptionKey(cardPublicKey, doorLoopSize)
	doorEncryptionKey := calculateEncryptionKey(doorPublicKey, cardLoopSize)

	answer := int64(-1)
	if cardEncryptionKey == doorEncryptionKey {
		answer = cardEncryptionKey
	}

	channel <- Result{
		answer:   answer,
		duration: time.Since(start).Nanoseconds(),
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

func loadPuzzleInput() (cardPublicKey int64, doorPublicKey int64) {
	//cardPublicKey = 5764801
	//doorPublicKey =  17807724
	cardPublicKey = 9093927
	doorPublicKey = 11001876
	return
}
