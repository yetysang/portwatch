package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/monitor"
	"portwatch/internal/ports"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	t	fmt.Fprintf(os.Stderr, "portwatch: failed to load config: %v\n", err)
		os.Exit(1)
	}

	scanner := ports.NewScanner()
	alertHandler := alert.NewHandler(os.Stdout, cfg)
	mon := monitor.New(scanner, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Fprintf(os.Stdout, "portwatch: starting (interval=%s)\n", cfg.Interval)

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	// Run an initial scan immediately.
	if err := tick(mon, alertHandler); err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: scan error: %v\n", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := tick(mon, alertHandler); err != nil {
				fmt.Fprintf(os.Stderr, "portwatch: scan error: %v\n", err)
			}
		case <-ctx.Done():
			fmt.Fprintln(os.Stdout, "portwatch: shutting down")
			return
		}
	}
}

func tick(mon *monitor.Monitor, h *alert.Handler) error {
	changes, err := mon.Poll()
	if err != nil {
		return err
	}
	for _, c := range changes {
		h.Handle(c)
	}
	return nil
}
