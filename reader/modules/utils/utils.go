package utils

import (
	"log"
	"time"
)

func Check(e error) {
	if e != nil {
		// log.Println(e)
		log.Fatal(e)
	}
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %v", name, elapsed.Seconds())
}
