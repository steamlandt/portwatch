package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/config"
	"github.com/example/portwatch/internal/monitor"
	"github.com/example/portwatch/internal/schedule"
	"github.com/example/portwatch/internal/state"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Printf("using defaults: %v", err)
		cfg = config.Default()
	}

	st, err := state.New(cfg.StateFile)
	if err != nil {
		log.Fatalf("state: %v", err)
	}

	al, err := alert.New(os.Stdout)
	if err != nil {
		log.Fatalf("alert: %v", err)
	}

	mon, err := monitor.New(cfg, st, al)
	if err != nil {
		log.Fatalf("monitor: %v", err)
	}

	sched := schedule.New(cfg.Interval)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("portwatch started (interval=%s range=%s)", cfg.Interval, cfg.PortRange)

	go sched.Start(func() {
		if err := mon.Run(); err != nil {
			log.Printf("monitor error: %v", err)
		}
	})

	<-sig
	log.Println("shutting down")
	sched.Stop()
}
