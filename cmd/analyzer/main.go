package main

import (
	"fmt"
	"os"

	"github.com/SanujaRubasinghe/ssdindexer/internal/scanner"
	"github.com/SanujaRubasinghe/ssdindexer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	path := "/"

	if len(os.Args) > 1 {
		path = os.Args[1]
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
