package main

import (
	"cardWithWords/internal/pkg/storage"
	"cardWithWords/internal/pkg/telegram"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jessevdk/go-flags"
)

type Opts struct {
	Token string `long:"token" description:"telegram api token" required:"true"`
}

var (
	opts         Opts
	chanOsSignal = make(chan os.Signal)
	chanDone     = make(chan struct{}, 2)
	chanError    = make(chan error, 32)
	wg           = new(sync.WaitGroup)
)

func dispatcher() {
	signal.Notify(chanOsSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-chanOsSignal

	log.Printf("[dispatcher] got %v\n", sig)
	chanDone <- struct{}{}
	chanDone <- struct{}{}
}

func main() {
	log.Println("start")

	if _, err := flags.Parse(&opts); err != nil {
		log.Fatalf("Parse flags: %v\n", err)
	}

	wg.Add(2)
	words := storage.Init()
	tb := telegram.Init(words, wg, chanDone, chanError)

	go dispatcher()

	// error getter and printer
	go func(wg *sync.WaitGroup, cDone chan struct{}, cErr chan error) {
		for {
			select {
			case <-cDone:
				wg.Done()
				return
			case msg := <-cErr:
				fmt.Println(msg)
			}
		}
	}(wg, chanDone, chanError)

	if err := tb.PlayCards(opts.Token); err != nil {
		log.Fatalf("couldn't initialize telegram bot api: %v\n", err)
	}

	wg.Wait()
	log.Println("done")
}
