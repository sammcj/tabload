package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (t *TabLoad) handleDownload() {
	params := map[string]interface{}{
		"repo_id":     t.repoIDEntry.Text,
		"revision":    t.revisionEntry.Text,
		"repo_type":   strings.ToLower(t.repoTypeDropdown.Selected),
		"folder_name": t.folderNameEntry.Text,
		"token":       t.tokenEntry.Text,
		"include":     strings.Split(t.includeEntry.Text, ","),
		"exclude":     strings.Split(t.excludeEntry.Text, ","),
	}

	downloadPath, err := t.client.Download(params)
	if err != nil {
		fmt.Println("Error downloading:", err)
		return
	}

	fmt.Println("Download completed successfully. Path:", downloadPath)
}

func (t *TabLoad) handleCancelDownload() {
	err := t.client.CancelDownload()
	if err != nil {
		fmt.Println("Error cancelling download:", err)
		return
	}

	fmt.Println("Download canceled successfully")
}

func (t *TabLoad) refreshData() error {
	var err error

	models, err := t.client.FetchModels()
	if err != nil {
		return fmt.Errorf("fetching models: %w", err)
	}
	t.modelsDropdown.Options = models

	loras, err := t.client.FetchLoras()
	if err != nil {
		return fmt.Errorf("fetching LoRAs: %w", err)
	}
	t.lorasDropdown.Options = loras

	t.refreshCurrentModel()
	t.refreshCurrentLoras()

	return nil
}

func (t *TabLoad) buildHFDownloaderTab() fyne.CanvasObject {
	t.repoIDEntry = widget.NewEntry()
	t.repoIDEntry.SetPlaceHolder("Repo ID")

	t.revisionEntry = widget.NewEntry()
	t.revisionEntry.SetPlaceHolder("Revision/Branch")

	t.repoTypeDropdown = widget.NewSelect([]string{"Model", "LoRA"}, func(selected string) {})

	t.folderNameEntry = widget.NewEntry()
	t.folderNameEntry.SetPlaceHolder("Folder Name")

	t.includeEntry = widget.NewEntry()
	t.includeEntry.SetPlaceHolder("Include Patterns")

	t.excludeEntry = widget.NewEntry()
	t.excludeEntry.SetPlaceHolder("Exclude Patterns")

	t.tokenEntry = widget.NewPasswordEntry()
	t.tokenEntry.SetPlaceHolder("HF Access Token")

	t.downloadButton = widget.NewButton("Download", t.handleDownload)
	t.cancelDownloadButton = widget.NewButton("Cancel Download", t.handleCancelDownload)

	return container.NewVBox(
		t.repoIDEntry,
		t.revisionEntry,
		t.repoTypeDropdown,
		t.folderNameEntry,
		t.includeEntry,
		t.excludeEntry,
		t.tokenEntry,
		t.downloadButton,
		t.cancelDownloadButton,
	)
}
