package randgen

import (
	"math"
	"math/rand"
	"time"
)

//this func generates numbers from 1 to 10 at intervals of 0.5 to 3 seconds and transmits them to the channel
func RandomGenerator(channelRandomNumber chan<- float64, closeChannel <-chan bool) {
	var result float64

	generationPeriodMsec := math.Round((rand.Float64()/1*(3.0-0.5) + 0.5) * 1000)
	generationPeriodSec := time.Duration(int64(generationPeriodMsec) * int64(time.Millisecond))
	timeGenerator := time.NewTicker(generationPeriodSec)
	for {
		select {
		case <-closeChannel:
			close(channelRandomNumber)
			timeGenerator.Stop()
			break
		case <-timeGenerator.C:
			rand.Seed(time.Now().UnixNano())
			generationPeriodMsec = math.Round((rand.Float64()/1*(3.0-0.5) + 0.5) * 1000)
			generationPeriodSec = time.Duration(int64(generationPeriodMsec) * int64(time.Millisecond))
			result = rand.Float64()/1*(10.0-1) + 1
			channelRandomNumber <- result
			timeGenerator.Reset(generationPeriodSec)
		}
	}
}
