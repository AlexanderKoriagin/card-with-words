package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"cardWithWords/internal/pkg/groq"
	"cardWithWords/internal/pkg/storage"
	"cardWithWords/internal/pkg/telegram"

	"github.com/jessevdk/go-flags"
)

type Opts struct {
	TokenTelegram string `long:"token-telegram" description:"telegram api token" required:"true"`
	TokenGroq     string `long:"token-groq" description:"groq api token" required:"true"`
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

	groqWords, err := groq.New(opts.TokenGroq)
	if err != nil {
		log.Fatalf("couldn't initialize groq client: %v\n", err)
	}

	localWords := storage.Init()
	tb := telegram.Init(
		localWords,
		groqWords,
		wg,
		chanDone,
		chanError,
	)

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

	err = tb.PlayCards(opts.TokenTelegram)
	if err != nil {
		log.Fatalf("couldn't initialize telegram bot api: %v\n", err)
	}

	wg.Wait()
	log.Println("done")
}
