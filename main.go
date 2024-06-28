package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/sammcj/tabload/logging"
	"github.com/sammcj/tabload/ui"
)

func main() {
	logging.Info("TabLoad started")

	a := app.NewWithID("com.sammcj.tabload")
	w := a.NewWindow("TabLoad")

	// Set initial window size
	w.Resize(fyne.Size{Width: 900, Height: 800})

	tabload := ui.NewTabLoad(w)

	// Build UI first
	tabload.BuildUI()

	// Perform auto-connect immediately after building the UI
	if tabload.ShouldAutoConnect() {
		tabload.AutoConnect()
	}

	// Set system tray menu
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("TabLoad",
			fyne.NewMenuItem("Open TabLoad", func() {
				logging.Info("Open TabLoad")
			}))
		desk.SetSystemTrayMenu(m)
	}

	w.ShowAndRun()

	logging.Info("TabLoad exited")
}
