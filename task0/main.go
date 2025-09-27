package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGUSR1)

	valueChannel := make(chan int64, 10)
	var values []int64

	go func() {
		for {
			s := <-signalChannel
			if s == syscall.SIGUSR1 {
				log.Println("Get SIGUSR")
			}
		}
	}()

	go func() {
		produceTicker := time.NewTicker(5 * time.Second)
		defer produceTicker.Stop()

		for {
			<-produceTicker.C
			randomValue := rand.Int63n(10)
			log.Printf("Get value: %d\n", randomValue)
			valueChannel <- randomValue
		}
	}()

	go func() {
		consumeTicker := time.NewTicker(30 * time.Second)
		defer consumeTicker.Stop()

		for {
			<-consumeTicker.C
			for {
				select {
				case n := <-valueChannel:
					values = append(values, n)
				default:
					goto FINISH
				}
			}
		FINISH:
			log.Printf("Processed elems: %v\n", values)
			values = values[:0]
		}
	}()

	select {}
}
