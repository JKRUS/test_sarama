package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jkrus/test_sarama/consumer/pkg/repository"
	"github.com/joho/godotenv"
	"github.com/kylelemons/go-gypsy/yaml"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	logrus.SetFormatter(new(logrus.JSONFormatter))
	f, err := os.OpenFile("logfile", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()
	logrus.SetOutput(f)

	config, err := readConfigFile("configs/config.yaml")
	if err != nil {
		logrus.Fatalf("error read config %s", err.Error())
	}
	db, err := initDB(*config, "configs/.env")
	if err != nil {
		logrus.Fatalf("error read config key %s", err.Error())
	}
	repos, err := repository.NewPostgresDB(*db)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	confConsumer, err := getConfigConsumer(*config)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	consumer, err := sarama.NewConsumer([]string{confConsumer.ip + ":" + confConsumer.port}, sarama.NewConfig())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			logrus.Fatalln(err)
		}
	}()
	//partitionConsumer, err := consumer.ConsumePartition("Topic_1", 0, sarama.OffsetNewest)
	partitionConsumer, err := consumer.ConsumePartition(confConsumer.topic, 0, -1)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	mess := repository.NewMsgPostgres(repos)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			_, err := mess.WriteMessage(msg.Value)
			if err != nil {
				logrus.Fatalln(err)
			}
			fmt.Printf("Consumed message offset %d\n", msg.Offset)
			fmt.Printf("Message:\t%s\n", string(msg.Value))
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}
	logrus.Printf("Consumed: %d\n", consumed)
}

func readConfigFile(path string) (*yaml.File, error) {
	config, err := yaml.ReadFile(path)
	if err != nil {
	}
	_, err = config.Get("kafka.ip")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("kafka.port")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("kafka.topic")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.host")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.port")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.username")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.dbname")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.sslmode")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func initDB(c yaml.File, path string) (*repository.Config, error) {
	db := &repository.Config{}
	var err error
	db.Host, err = c.Get("db.host")
	if err != nil {
		return nil, err
	}
	db.Port, err = c.Get("db.port")
	if err != nil {
		return nil, err
	}
	db.UserName, err = c.Get("db.username")
	if err != nil {
		return nil, err
	}
	db.DBName, err = c.Get("db.dbname")
	if err != nil {
		return nil, err
	}
	db.SSLMode, err = c.Get("db.sslmode")
	if err != nil {
		return nil, err
	}
	if err = godotenv.Load(path); err != nil {
		return nil, err
	}
	db.Password = os.Getenv("DB_PASSWORD")
	return db, nil
}

type configConsumer struct {
	ip    string
	port  string
	topic string
}

func getConfigConsumer(c yaml.File) (*configConsumer, error) {
	conf := &configConsumer{}
	var err error
	conf.ip, err = c.Get("kafka.ip")
	if err != nil {
		return nil, err
	}
	conf.port, err = c.Get("kafka.port")
	if err != nil {
		return nil, err
	}
	conf.topic, err = c.Get("kafka.topic")
	if err != nil {
		return nil, err
	}
	return conf, nil
}
