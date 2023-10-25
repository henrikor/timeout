package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/gookit/color"
	"github.com/gosuri/uiprogress"
	"github.com/hajimehoshi/go-mp3"
)

func main() {
	args := os.Args
	a := os.Args[1]

	// Hvis ingen argumenter er gitt, skriv ut en melding
	if len(args) == 1 {
		fmt.Println("Ingen argumenter ble gitt.")
		os.Exit(1)
	}

	fmt.Println("Arg1: " + a)

	var sec int
	konvfeil := "Feil ved konvertering: "

	if strings.Contains(a, ":") {
		parts := strings.Split(a, ":")
		m := parts[0]
		s := parts[1]

		fmt.Println("min: " + m + " sec: " + s)
		min, err := strconv.Atoi(m)
		if err != nil {
			log.Fatal(konvfeil, err)
		}
		se, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(konvfeil, err)
		}

		sec = min*60 + se
	} else {
		var err error
		sec, err = strconv.Atoi(a)
		if err != nil {
			log.Fatal(konvfeil, err)
		}

	}

	fmt.Println("Sec totalt: ", sec)

	if sec < 30 {
		color.Errorln("Number of seconds must be more then 30 seconds, I got only %d seconds", sec)
		os.Exit(1)
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	go doTheSounds(sec)

	barwait := sec + 7
	// if sec < 40 {
	// } else {
	// 	barwait = sec
	// }

	// //////////////////////////////
	uiprogress.Start()                // start rendering
	bar := uiprogress.AddBar(barwait) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	bar.AppendCompleted()
	bar.PrependElapsed()

	for bar.Incr() {
		time.Sleep(time.Second * 1)
	} ////////////////////////////////

}

func doTheSounds(sec int) {
	otoCtx, readyChan, err := createContext()

	<-readyChan

	x := sec - 30
	biiip := "./beep-01a.mp3"
	cow := "./cow.mp3"

	time.Sleep(time.Duration(x) * time.Second)
	playFile(biiip, otoCtx, err)
	time.Sleep(26 * time.Second)
	playFile(biiip, otoCtx, err)
	time.Sleep(1 * time.Second)
	playFile(biiip, otoCtx, err)
	time.Sleep(1 * time.Second)
	playFile(biiip, otoCtx, err)
	time.Sleep(2 * time.Second)
	playFile(cow, otoCtx, err)
}

func playFile(file string, otoCtx *oto.Context, err error) {
	decodedMp3 := mkDecodeMp3(file)
	player := otoCtx.NewPlayer(decodedMp3)
	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}

func createContext() (*oto.Context, chan struct{}, error) {
	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	return otoCtx, readyChan, err
}

func mkDecodeMp3(file string) *mp3.Decoder {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		txt := fmt.Sprintf("reading %s failed", file+err.Error())
		panic(txt)
	}

	fileBytesReader := bytes.NewReader(fileBytes)
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}
	return decodedMp3
}
