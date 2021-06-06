package main

import (
	"github.com/jkrus/test_sarama/app/pkg/controller"
	"github.com/jkrus/test_sarama/app/pkg/repository"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

func routes(r *httprouter.Router) {
	r.ServeFiles("/../public/*filepath", http.Dir("/../public"))
	r.GET("/", controller.StartPage)
	r.GET("/data", controller.GetTable)
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	logrus.SetFormatter(new(logrus.JSONFormatter))
	f, err := os.OpenFile("logfile", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()
	logrus.SetOutput(f)

	err = repository.NewPostgresDB()
	if err != nil {
		logrus.Fatal(err)
	}

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	r := httprouter.New()
	routes(r)

	err = http.ListenAndServe("localhost:8000", r)
	if err != nil {
		logrus.Fatal(err)
	}
}
