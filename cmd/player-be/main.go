package main

import (
	"player-be/internal/boot"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	err := boot.HTTP()
	if err != nil {
		log.Fatalln(err.Error())
	}

}
