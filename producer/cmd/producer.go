package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jkrus/test_sarama/producer/pkg/parsers/config_yaml"
	"github.com/jkrus/test_sarama/producer/pkg/parsers/table"
	"github.com/jkrus/test_sarama/producer/pkg/randgen"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type parameters struct {
	generate bool
	file     bool
	numbers  int
}

func (p *parameters) registerFlags() {
	flag.BoolVar(&p.generate, "g", false, "generate random numbers")
	flag.BoolVar(&p.file, "f", false, "read data from file")
	flag.IntVar(&p.numbers, "n", 16, "the number of iterations for a generator or file pass ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [dir]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}

func selectData(p parameters, closeChannel <-chan syscall.Signal, channelData chan<- []byte) {
	if p.generate && p.file {
		//logrus.Fatalf("error %s", err.Error())
	}
	//p.generate = true
	if p.generate {
		err := generateData(p.numbers, closeChannel, channelData)
		if err != nil {
			logrus.Fatalf("error %s", err.Error())
		}
	}
	//p.file = true
	if p.file {
		err := readData(p.numbers, closeChannel, channelData)
		if err != nil {
			logrus.Fatalf("error %s", err.Error())
		}
	}
}

func generateData(n int, closeChannel <-chan syscall.Signal, channelData chan<- []byte) error {
	channelRandomGenerator := make(chan int)
	channelCloseRandomGenerator := make(chan bool)
	go randgen.RandomGenerator(n, channelRandomGenerator, channelCloseRandomGenerator)
	i := 0
loop:
	for {
		select {
		case <-closeChannel:
			channelCloseRandomGenerator <- true
			fmt.Println("\nГенерация данных окончена")
		case randomNumber := <-channelRandomGenerator:
			tabRes := make(map[string][]int)
			tabRes["0"] = append(tabRes["0"], randomNumber)
			jitem, err := json.Marshal(tabRes)
			if err != nil {
				fmt.Println(err.Error())
			}
			channelData <- jitem
			i++
		}
		if i == n {
			break loop
		}
	}
	return nil
}

func readData(n int, closeChannel <-chan syscall.Signal, channelData chan<- []byte) error {
	i := 0
LoopRead:
	for {
		select {
		case <-closeChannel:
			fmt.Println("\nЧтение данных окончено")
			break
		default:
			i++
			t := table.NewTable()
			tab, err := t.ReadTable("data.txt")
			if err != nil {
				logrus.Fatalf("error %s", err.Error())
			}
			tabRes := make(map[string][]int)
			for k := 0; k < len(tab); k++ {
				s := strconv.Itoa(k)
				tabRes[s] = tab[k]
			}
			jitem, err := json.Marshal(tabRes)
			if err != nil {
				fmt.Println(err.Error())
			}
			channelData <- jitem
			if i == n {
				break LoopRead
			}
		}
	}
	return nil
}

func main() {
	// Trap SIGINT to trigger a shutdown.
	channelEndProgram := make(chan os.Signal, 1)
	channelCloseGoroutines := make(chan syscall.Signal)
	signal.Notify(channelEndProgram, os.Interrupt)
	channelData := make(chan []byte, 5)
	params := parameters{}
	params.registerFlags()
	flag.Parse()

	c := config_yaml.NewConfig()
	conf, err := c.ReadFile("configs/config.yaml")
	if err != nil {
		logrus.Fatalf("error %s", err.Error())
	}

	go selectData(params, channelCloseGoroutines, channelData)

	producer, err := sarama.NewAsyncProducer([]string{conf[0]}, nil)
	if err != nil {
		logrus.Fatalf("error connect to Kafka, %s", err.Error())
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logrus.Fatalf("error close producer, %s", err.Error())
		}
	}()

	var enqueued, producerErrors int
mainLoop:
	for {
		select {
		case data := <-channelData:
			producer.Input() <- &sarama.ProducerMessage{Topic: conf[1], Key: nil, Value: sarama.StringEncoder(data)}
			fmt.Println(sarama.StringEncoder(data))
			enqueued++
		case err := <-producer.Errors():
			logrus.Println("Failed to produce message", err)
			producerErrors++
		case <-channelEndProgram:
			channelCloseGoroutines <- unix.SIGQUIT
			close(channelData)
			fmt.Println("\nПользователь завершил выполнение программы")
		}
		if enqueued == params.numbers {
			close(channelData)
			logrus.Printf("Enqueued: %d; errors: %d\n", enqueued, producerErrors)
			break mainLoop
		}
	}
}
