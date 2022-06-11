package main

import (
	"log"
	"os"
	"fmt"
	"github.com/faiface/beep"
)

func main() {
	f, err := os.Open("cow.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("test")
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
}
