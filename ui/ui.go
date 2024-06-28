package ui

import (
	"errors"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sammcj/tabload/api"
	"github.com/sammcj/tabload/logging"
)

func NewTabLoad(w fyne.Window) *TabLoad {
	initConfig()
	t := &TabLoad{window: w}

	// initialise UI elements
	t.apiURLEntry = widget.NewEntry()
	t.adminKeyEntry = widget.NewPasswordEntry()
	t.connectButton = widget.NewButton("Connect", t.handleConnect)

	// initialise other necessary UI elements
	t.modelsDropdown = widget.NewSelect([]string{}, func(selected string) {})
	t.currentModelLabel = widget.NewLabel("")
	t.currentLorasLabel = widget.NewLabel("")

	// initialise all entry fields and checkboxes
	t.maxSeqLenEntry = widget.NewEntry()
	t.maxSeqLenCheck = widget.NewCheck("", nil)
	t.overrideBaseSeqLenEntry = widget.NewEntry()
	t.overrideBaseSeqLenCheck = widget.NewCheck("", nil)
	t.cacheSizeEntry = widget.NewEntry()
	t.cacheSizeCheck = widget.NewCheck("", nil)
	t.gpuSplitEntry = widget.NewEntry()
	t.gpuSplitCheck = widget.NewCheck("", nil)
	t.gpuSplitAutoCheck = widget.NewCheck("", nil)
	t.ropeScaleEntry = widget.NewEntry()
	t.ropeScaleCheck = widget.NewCheck("", nil)
	t.ropeAlphaEntry = widget.NewEntry()
	t.ropeAlphaCheck = widget.NewCheck("", nil)
	t.cacheModeDropdown = widget.NewSelect([]string{"Q4", "Q6", "Q8", "FP16"}, nil)
	t.promptTemplateEntry = widget.NewEntry()
	t.promptTemplateCheck = widget.NewCheck("", nil)
	t.numExpertsPerTokenEntry = widget.NewEntry()
	t.numExpertsPerTokenCheck = widget.NewCheck("", nil)
	t.draftModelNameEntry = widget.NewEntry()
	t.draftModelNameCheck = widget.NewCheck("", nil)
	t.draftRopeScaleEntry = widget.NewEntry()
	t.draftRopeScaleCheck = widget.NewCheck("", nil)
	t.draftRopeAlphaEntry = widget.NewEntry()
	t.draftRopeAlphaCheck = widget.NewCheck("", nil)
	t.draftCacheModeDropdown = widget.NewSelect([]string{"Q4", "Q6", "Q8", "FP16"}, nil)
	t.fasttensorsCheck = widget.NewCheck("", nil)
	t.autosplitReserveEntry = widget.NewEntry()
	t.autosplitReserveCheck = widget.NewCheck("", nil)
	t.chunkSizeEntry = widget.NewEntry()
	t.chunkSizeCheck = widget.NewCheck("", nil)

	// Set the API URL from config
	config := t.GetConfig()

	// Set the API URL from config
	serverURL := config.LastConnectedServer
	if serverURL == "" {
		serverURL = config.APIURL
	}
	if serverURL == "" {
		serverURL = "http://localhost:5000" // Default URL only if no server is set
	}
	logging.Info(fmt.Sprintf("Using server URL: %s", serverURL))
	t.apiURLEntry.SetText(serverURL)

	// initialise the client with the server URL
	t.client = api.NewClient(serverURL, "")

	// Load default parameters
	t.LoadDefaultParams()

	return t
}

func (t *TabLoad) BuildUI() {
	logging.Debug("Building UI")
	if t.client == nil {
		logging.Error("Client is not initialised", nil)
		dialog.ShowError(errors.New("client is not initialised"), t.window)
		return
	}

	tabs := container.NewAppTabs(
		container.NewTabItem("Connection", t.buildConnectionTab()),
		container.NewTabItem("Model", t.buildModelTab()),
		container.NewTabItem("LoRAs", t.buildLorasTab()),
		container.NewTabItem("HF Downloader", t.buildHFDownloaderTab()),
		container.NewTabItem("Presets", t.buildPresetTab()),
		container.NewTabItem("Settings", t.buildSettingsTab()),
		container.NewTabItem("Advanced", t.buildAdvancedSettingsTab()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	saveDefaultsButton := widget.NewButton("Save as Default", func() {
		if err := t.SaveDefaultParams(); err != nil {
			dialog.ShowError(err, t.window)
		} else {
			dialog.ShowInformation("Success", "Parameters saved as default", t.window)
		}
	})

	toolbar := t.createToolbar(saveDefaultsButton)

	t.currentModelInfo = container.NewVBox(widget.NewLabel("No model loaded"))

	statusBar := t.createStatusBar()

	// Create a compact horizontal container for model info
	modelInfoBar := container.NewHBox()
	t.currentModelInfo = modelInfoBar

	mainContent := container.NewBorder(nil, nil, nil, nil, tabs)

	t.logsPane = container.NewBorder(
		widget.NewLabel("Logs"),
		nil, nil, nil,
		container.NewScroll(widget.NewMultiLineEntry()),
	)
	// hide the logs pane by default
	t.logsPane.Hide()

	splitContainer := container.NewHSplit(mainContent, t.logsPane)
	splitContainer.Offset = 1 // Hide logs pane by default

	content := container.NewBorder(
		toolbar,
		container.NewVBox(modelInfoBar, statusBar),
		nil,
		nil,
		splitContainer,
	)

	// disable scrollbars
	t.window.SetPadded(false)

	t.window.SetContent(content)
	t.window.Resize(fyne.NewSize(1024, 768))
	t.SetInitialFocus()

	t.ready = true
	logging.Info("UI built successfully")

	// Attempt auto-connect after UI is built
	if config.AutoConnect && config.LastConnectedServer != "" {
		go t.AutoConnect()
	}
}
func (t *TabLoad) createToolbar(...*widget.Button) *widget.Toolbar {

	return widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			logging.Debug("Creating new model")
			dialog.ShowInformation("New Model", "Create a new model configuration", t.window)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			logging.Debug("Saving preset")
			t.handleSavePreset()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			logging.Debug("Downloading model")
			t.handleLoadModel()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			logging.Debug("Opening settings dialog")
			t.showSettingsDialog()
		}),
		widget.NewToolbarAction(theme.DocumentIcon(), func() {
			logging.Debug("Toggling logs pane")
			t.toggleLogsPane()
		}),
	)
}

func (t *TabLoad) createStatusBar() fyne.CanvasObject {
	logging.Debug("Creating status bar")
	connectionStatus := widget.NewLabel("Not connected")
	t.connectionStatus = connectionStatus

	return container.NewHBox(
		widget.NewLabel("Status:"),
		connectionStatus,
	)
}

func (t *TabLoad) SetInitialFocus() {
	logging.Debug("Setting initial focus")
	lastServer, err := t.loadLastConnectedServer()
	if t.window.Canvas().Focused() == nil {
		if err == nil && lastServer != "" {
			t.window.Canvas().Focus(t.connectButton)
		} else {
			t.window.Canvas().Focus(t.apiURLEntry)
		}
	} else {
		t.window.Canvas().Focus(nil)
	}
}

func (t *TabLoad) RefreshUI() {
	logging.Debug("Refreshing UI")
	t.modelsDropdown.Refresh()
	t.lorasDropdown.Refresh()
	t.currentModelLabel.Refresh()
	t.currentLorasLabel.Refresh()
	t.refreshCurrentModel()
}

func (t *TabLoad) toggleLogsPane() {
	if t.logsTextArea == nil {
		t.logsTextArea = widget.NewMultiLineEntry()
		t.logsTextArea.SetText("Logs will appear here...")
		t.logsTextArea.Disable()
	}

	if t.logsPane == nil {
		t.logsPane = container.NewBorder(
			widget.NewLabel("Logs"),
			nil, nil, nil,
			container.NewScroll(t.logsTextArea),
		)
	}

	t.showingLogs = !t.showingLogs
	t.refreshMainLayout()

	if t.showingLogs {
		go func() {
			logStream := logging.StreamLogs()
			for log := range logStream {
				t.logsTextArea.SetText(t.logsTextArea.Text + "\n" + log)
				t.logsTextArea.CursorRow = len(strings.Split(t.logsTextArea.Text, "\n")) - 1
			}
		}()
	}
}
func (t *TabLoad) refreshMainLayout() {
	mainContent := t.window.Content().(*fyne.Container)
	splitContainer := mainContent.Objects[0].(*container.Split)
	splitContainer.Offset = 0.7 // Show logs pane (70% main content, 30% logs)

	if t.showingLogs {
		t.logsPane.Show()
	} else {
		splitContainer.Offset = 1 // Hide logs pane
		t.logsPane.Hide()
	}

	mainContent.Refresh()
}
