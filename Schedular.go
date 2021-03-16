package main

import (
	"fmt"
	"github.com/robfig/cron"
)

func StartScheduler()  {
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		fmt.Println("Updating Recogniser ... ")
		UpdateRecogniser()
	})
	c.Start()
}
