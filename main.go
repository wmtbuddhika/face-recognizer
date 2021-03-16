package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	dateFormat = "2006_01_02"
)

func init() {
	var err error
	err = godotenv.Load() // can pass the location. default root
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetFlags(log.Ldate | log.Ltime)
	file, err := os.OpenFile("logs/"+time.Now().Format(dateFormat)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err == nil {
		log.SetOutput(file)
	}
}

func main() {
	fmt.Println("Starting Scheduler ... ")
	StartScheduler()

	fmt.Println("Starting Recogniser ... ")
	StartRecognition()
}