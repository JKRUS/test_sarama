package randgen

import (
	"math"
	"math/rand"
	"time"
)

//this func generates numbers from 1 to 10 at intervals of 0.5 to 3 seconds and transmits them to the channel
func RandomGenerator(number int, channelRandomNumber chan<- int, closeChannel <-chan bool) {
	var result, i int

	generationPeriodMsec := math.Round((rand.Float64()/1*(3.0-0.5) + 0.5) * 1000)
	generationPeriodSec := time.Duration(int64(generationPeriodMsec) * int64(time.Millisecond))
	timeGenerator := time.NewTicker(generationPeriodSec)
randLoop:
	for {
		select {
		case <-closeChannel:
			timeGenerator.Stop()
			break
		case <-timeGenerator.C:
			rand.Seed(time.Now().UnixNano())
			generationPeriodMsec = math.Round((rand.Float64()/1*(3.0-0.5) + 0.5) * 1000)
			generationPeriodSec = time.Duration(int64(generationPeriodMsec) * int64(time.Millisecond))
			result = int((rand.Float64()/1*(10.0-1)+1)*100 + 0.5)
			channelRandomNumber <- result
			timeGenerator.Reset(generationPeriodSec)
			i++
			if i == number {
				timeGenerator.Stop()
				break randLoop
			}
		}
	}
}
