package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/JKRUS/test_sarama/producer/pkg/parsers/config_yaml"
	"github.com/JKRUS/test_sarama/producer/pkg/parsers/table"
	"github.com/JKRUS/test_sarama/producer/pkg/randgen"
	"github.com/Shopify/sarama"
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
	flag.IntVar(&p.numbers, "n", 5, "the number of iterations for a generator or file pass ")

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
		err := generateData(p, closeChannel, channelData)
		if err != nil {
			logrus.Fatalf("error %s", err.Error())
		}
	}
	//p.file = true
	if p.file {
		err := readData(p, channelData)
		if err != nil {
			logrus.Fatalf("error %s", err.Error())
		}
	}
}

func generateData(p parameters, closeChannel <-chan syscall.Signal, channelData chan<- []byte) error {
	channelRandomGenerator := make(chan float64)
	channelCloseRandomGenerator := make(chan bool)
	go randgen.RandomGenerator(channelRandomGenerator, channelCloseRandomGenerator)
	i := 0
loop:
	for {
		i++
		select {
		case randomNumber, ok := <-channelRandomGenerator:
			if !ok {
				fmt.Println("Канал генерации и передачи данных закрыт")
				close(channelData)
				break loop
			} else {
				tabRes := make(map[string]string)
				s := fmt.Sprintf("%.2f", randomNumber)
				tabRes["0"] = s
				jitem, err := json.Marshal(tabRes)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(s)
				channelData <- jitem
				if i == p.numbers {
					fmt.Println("\nГенерация данных окончена")
					channelCloseRandomGenerator <- true
				}
			}
		case <-closeChannel:
			channelCloseRandomGenerator <- true
		}
	}
	return nil
}

func readData(p parameters, channelData chan<- []byte) error {
	i := 0
	var s string
	for {
		i++
		t := table.NewTable()
		tab, err := t.ReadTable("data.txt")
		if err != nil {
			logrus.Fatalf("error %s", err.Error())
		}
		tabRes := make(map[string]string)
		for k := 0; k < len(tab); k++ {
			s := strconv.Itoa(k)
			tabRes[s] = tab[k]
		}
		jitem, err := json.Marshal(tabRes)
		if err != nil {
			fmt.Println(err.Error())
		}
		for k := 0; k < len(tab); k++ {
			s = fmt.Sprintf("%s", tab[k])
			channelData <- jitem
			fmt.Println(s)
		}
		if i == p.numbers {
			fmt.Println("\nЧтение данных окончено")
			close(channelData)
			break
		}
	}
	return nil
}

func main() {
	// Trap SIGINT to trigger a shutdown.
	channelEndProgram := make(chan os.Signal, 1)
	channelCloseGoroutines := make(chan syscall.Signal)
	signal.Notify(channelEndProgram, os.Interrupt)

	channelData := make(chan []byte)

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
	for {

		select {
		case data, ok := <-channelData:
			if !ok {
				logrus.Printf("Enqueued: %d; errors: %d\n", enqueued, producerErrors)
				os.Exit(123)
			}
			producer.Input() <- &sarama.ProducerMessage{Topic: conf[1], Key: nil, Value: sarama.StringEncoder(data)}
			enqueued++
		case err := <-producer.Errors():
			logrus.Println("Failed to produce message", err)
			producerErrors++

		case <-channelEndProgram:
			channelCloseGoroutines <- unix.SIGQUIT
			fmt.Println("\nПользователь завершил выполнение программы")
		}
	}
}
