package main

import (
	"github.com/JKRUS/test_sarama/consumer"
	"github.com/sirupsen/logrus"
)

func main() {
	srv := new(consumer.Server)
	if err := srv.Run("8080"); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())

	}
}
