package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func startHealthCheck() {
	// Check each server if it's alive every 5 seconds
	scheduler := gocron.NewScheduler(time.Local)

	for _, backend := range backendList {
		_, err := scheduler.Every(5).Seconds().Do(func(server *server) {
			alive := server.isAlive()

			if alive {
				log.Printf("%s lives.", server.Name)
			} else {
				log.Printf("%s is dead.", server.Name)
			}
		}, backend)
		if err != nil {
			log.Fatalln(err)
		}
	}
	scheduler.StartAsync()
}
