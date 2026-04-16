package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func main() {
	cfgPath := flag.String("config", "", "path to JSON config file (optional)")
	flag.Parse()

	var cfg *config.Config
	var err error
	if *cfgPath != "" {
		cfg, err = config.Load(*cfgPath)
		if err != nil {
			log.Fatalf("portwatch: failed to load config: %v", err)
		}
	} else {
		cfg = config.Default()
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("portwatch: invalid config: %v", err)
	}

	var alertOut = os.Stdout
	if cfg.AlertOutput != "" {
		f, err := os.OpenFile(cfg.AlertOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("portwatch: cannot open alert output: %v", err)
		}
		defer f.Close()
		alertOut = f
	}

	sc := scanner.New(cfg.PortRange.Start, cfg.PortRange.End)
	st, err := state.New(cfg.StatePath)
	if err != nil {
		log.Fatalf("portwatch: state init failed: %v", err)
	}
	al := alert.New(alertOut)
	mon := monitor.New(sc, st, al)

	fmt.Printf("portwatch: starting — scanning ports %d-%d every %v\n",
		cfg.PortRange.Start, cfg.PortRange.End, cfg.Interval)

	mon.Run(cfg.Interval)
}
