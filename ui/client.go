package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sammcj/tabload/api"
	"github.com/sammcj/tabload/logging"
)

func (t *TabLoad) SetClient(client *api.Client) {
	t.client = client
	if t.ready && t.presetDropdown != nil {
		t.refreshPresetList()
	}
}
func (t *TabLoad) handleConnect() {
	url := t.apiURLEntry.Text
	key := t.adminKeyEntry.Text

	t.client.BaseURL = url
	t.client.AdminKey = key

	go func() {
		err := t.refreshData()
		if err != nil {
			logging.Error(fmt.Sprintf("Error refreshing data for %s", url), err)
			func() {
				dialog.ShowError(err, t.window)
				t.connectionStatus.SetText("Connection failed")
			}()
			return
		}

		// Save the last connected server
		if err := t.saveLastConnectedServer(url); err != nil {
			logging.Error("Failed to save last connected server", err)
		}

		func() {
			t.RefreshUI()
			t.window.SetTitle("TabLoad (connected to " + url + ")")
			t.connectButton.Hide()
			t.connectionStatus.SetText("Connected to " + url)
		}()
	}()
}

func (t *TabLoad) buildConnectionTab() fyne.CanvasObject {
	t.apiURLEntry = widget.NewEntry()
	t.apiURLEntry.SetPlaceHolder("TabbyAPI Endpoint URL")

	t.adminKeyEntry = widget.NewPasswordEntry()
	t.adminKeyEntry.SetPlaceHolder("Admin Key")

	t.connectButton = widget.NewButton("Connect", t.handleConnect)

	lastServer, err := t.loadLastConnectedServer()
	if err == nil && lastServer != "" {
		t.apiURLEntry.SetText(lastServer)
	}

	return container.NewVBox(
		t.apiURLEntry,
		t.adminKeyEntry,
		t.connectButton,
	)
}

func (t *TabLoad) AutoConnect() {
	if !t.ready {
		logging.Warn("Cannot auto-connect: UI is not fully initialised")
		return
	}

	lastServer := config.LastConnectedServer
	logging.Info(fmt.Sprintf("Auto-connect attempt. Last server: %s", lastServer))

	if lastServer != "" {
		logging.Info(fmt.Sprintf("Auto-connecting to last connected server: %s", lastServer))
		t.apiURLEntry.SetText(lastServer)
		t.client.BaseURL = lastServer
		t.handleConnect()
	} else {
		logging.Warn("Auto-connect failed: No last connected server found")
	}
}
