package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func startHealthCheck(numberOfSeconds int) {
	currentTime := time.Now()

	var backendList []*server

	// Health check first server list
	// if day in the month is < 15
	if currentTime.Day() < 15 {
		backendList = serverListFirst
	} else {
		backendList = serverListSecond
	}

	// Check each server if it's alive every n seconds
	scheduler := gocron.NewScheduler(time.Local)

	for _, backend := range backendList {
		_, err := scheduler.Every(numberOfSeconds).Seconds().Do(func(server *server) {
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
