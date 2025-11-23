package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SanujaRubasinghe/ssdindexer/internal/scanner"
	"github.com/SanujaRubasinghe/ssdindexer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Default to current working directory if no path is provided
	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Use first argument as path if provided
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// Convert to absolute path
	path, err = filepath.Abs(path)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		os.Exit(1)
	}

	// Create the TUI model
	m := ui.New()

	// Channel for live updates from scanner
	updateChan := make(chan scanner.FileStats)
	doneChan := make(chan struct{})

	// Start scanning in a goroutine
	go func() {
		stats, err := scanner.Scan(path, updateChan) // new ScanLive function
		if err != nil {
			fmt.Println("Error:", err)
			close(doneChan)
			return
		}
		// final stats
		updateChan <- stats
		close(doneChan)
	}()

	// Bubble Tea program
	prog := tea.NewProgram(m)

	// Start a goroutine to feed updates to TUI
	go func() {
		for {
			select {
			case stats, ok := <-updateChan:
				if !ok {
					return
				}
				prog.Send(ui.UpdateMsg(stats))
			case <-doneChan:
				prog.Send(scanner.DoneMsg{})
			}
		}
	}()

	if err := prog.Start(); err != nil {
		panic(err)
	}
}
