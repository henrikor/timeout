package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/gookit/color"
	"github.com/hajimehoshi/go-mp3"
	"golang.org/x/crypto/ssh/terminal"

	"atomicgo.dev/cursor"
	"github.com/hedzr/progressbar"
)

var whichStepper = 1
var count = 0

func forAllSteppers(barwait int) {
	tasks := progressbar.NewTasks(progressbar.New())
	defer tasks.Close()

	max := count
	_, h, _ := terminal.GetSize(int(os.Stdout.Fd()))
	if max >= h {
		max = h
	}

	for i := whichStepper; i < whichStepper+max; i++ {
		tasks.Add(
			progressbar.WithTaskAddBarOptions(
				progressbar.WithBarStepper(i),
				progressbar.WithBarUpperBound(100),
				progressbar.WithBarWidth(64),
				// progressbar.WithBarTextSchema(schema),
			),
			progressbar.WithTaskAddBarTitle(string(strconv.AppendInt([]byte("Task "), int64(i), 10))), // fmt.Sprintf("Task %v", i)),
			progressbar.WithTaskAddOnTaskProgressing(func(bar progressbar.PB, exitCh <-chan struct{}) {
				for ub, ix := bar.UpperBound(), int64(0); ix < ub; ix++ {
					ms := time.Duration(barwait) //nolint:gosec //just a demo
					time.Sleep(time.Millisecond * ms * 10)
					bar.Step(1)
				}
			}),
		)
	}

	tasks.Wait() // start waiting for all tasks completed gracefully
}

func doTheSounds(sec int) {
	otoCtx, readyChan, err := createContext()

	<-readyChan

	x := sec - 30

	path := "./"
	if os.Getenv("TIMEOUT_PATH") != "" {
		path = os.Getenv("TIMEOUT_PATH")
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	biiip := path + "beep-01a.mp3"
	cow := path + "cow.mp3"

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
func main() {
	args := os.Args

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("Got signal: ", sig)
			cursor.Show()
			os.Exit(0)
		}
	}()
	// Hvis ingen argumenter er gitt, skriv ut en melding
	if len(args) == 1 {
		fmt.Println("Ingen argumenter ble gitt.")
		fmt.Println("Oppgi antall sekunder, eller minutter:sekunder")
		fmt.Println("Eks:")
		fmt.Println("  timeout 30")
		fmt.Println("  timeout 2:00")
		os.Exit(1)
	}
	a := os.Args[1]

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

	cursor.Hide()
	defer cursor.Show()

	count = 1
	whichStepper = 1
	forAllSteppers(barwait)
}
