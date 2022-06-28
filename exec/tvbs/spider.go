package main

import (
	"io"
	"log"
	"os"
	"github.com/sirupsen/logrus"

	mongodb "annstudio.com/spiders/crawler/db" 
	tvbs "annstudio.com/spiders/crawler/tvbs"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{}) //   .JSONFormatter{}
	logrus.SetLevel(logrus.DebugLevel)

	writer1 := os.Stdout
	writer2, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	logrus.SetOutput(io.MultiWriter(writer1, writer2))

	data_list := tvbs.Get_list(10) // 取前幾個標題文章
	logrus.Info("result amount:", len(data_list))
	mongodb.Upsert_data(data_list)
}
