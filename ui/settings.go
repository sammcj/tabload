package ui

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/viper"
)

func (t *TabLoad) buildSettingsTab() fyne.CanvasObject {
	autoConnectCheck := widget.NewCheck("Auto-connect on startup", func(checked bool) {
		t.saveAutoConnectSetting(checked)
	})
	return container.NewVBox(autoConnectCheck)
}

func (t *TabLoad) ShouldAutoConnect() bool {
	return viper.GetBool("auto_connect")
}

func (t *TabLoad) buildAdvancedSettingsTab() fyne.CanvasObject {
	// Sampling parameters
	temperatureSlider := widget.NewSlider(0, 2)
	temperatureSlider.SetValue(1)

	topKEntry := widget.NewEntry()
	topKEntry.SetPlaceHolder("Top K (e.g., 40)")

	topPSlider := widget.NewSlider(0, 1)
	topPSlider.SetValue(1)

	minPSlider := widget.NewSlider(0, 1)
	minPSlider.SetValue(0)

	topASlider := widget.NewSlider(0, 1)
	topASlider.SetValue(0)

	tfsSlider := widget.NewSlider(0, 1)
	tfsSlider.SetValue(1)

	typicalPSlider := widget.NewSlider(0, 1)
	typicalPSlider.SetValue(1)

	repetitionPenaltySlider := widget.NewSlider(1, 2)
	repetitionPenaltySlider.SetValue(1)

	presencePenaltySlider := widget.NewSlider(-2, 2)
	presencePenaltySlider.SetValue(0)

	frequencyPenaltySlider := widget.NewSlider(-2, 2)
	frequencyPenaltySlider.SetValue(0)

	mirostatModeSelect := widget.NewSelect([]string{"0", "1", "2"}, func(s string) {})
	mirostatModeSelect.SetSelected("0")

	mirostatTauSlider := widget.NewSlider(0, 10)
	mirostatTauSlider.SetValue(5)

	mirostatEtaSlider := widget.NewSlider(0, 1)
	mirostatEtaSlider.SetValue(0.1)

	// Other settings
	streamingCheck := widget.NewCheck("Enable Streaming", func(bool) {})

	grammarEntry := widget.NewMultiLineEntry()
	grammarEntry.SetPlaceHolder("Enter grammar string here")

	logitBiasEntry := widget.NewEntry()
	logitBiasEntry.SetPlaceHolder("Logit bias (format: 'tokenID:bias,tokenID:bias')")

	negativePromptEntry := widget.NewMultiLineEntry()
	negativePromptEntry.SetPlaceHolder("Enter negative prompt here")

	jsonModeCheck := widget.NewCheck("Enable JSON Mode", func(bool) {})

	// Speculative decoding
	speculativeNgramCheck := widget.NewCheck("Enable Speculative Decoding", func(bool) {})

	saveButton := widget.NewButton("Save Advanced Settings", func() {
		// Handle saving advanced settings
		t.saveAdvancedSettings()
	})

	// Create a grid layout for sampling parameters
	samplingGrid := container.NewGridWithColumns(2,
		widget.NewLabel("Temperature:"), temperatureSlider,
		widget.NewLabel("Top K:"), topKEntry,
		widget.NewLabel("Top P:"), topPSlider,
		widget.NewLabel("Min P:"), minPSlider,
		widget.NewLabel("Top A:"), topASlider,
		widget.NewLabel("TFS:"), tfsSlider,
		widget.NewLabel("Typical P:"), typicalPSlider,
		widget.NewLabel("Repetition Penalty:"), repetitionPenaltySlider,
		widget.NewLabel("Presence Penalty:"), presencePenaltySlider,
		widget.NewLabel("Frequency Penalty:"), frequencyPenaltySlider,
		widget.NewLabel("Mirostat Mode:"), mirostatModeSelect,
		widget.NewLabel("Mirostat Tau:"), mirostatTauSlider,
		widget.NewLabel("Mirostat Eta:"), mirostatEtaSlider,
	)

	// Combine all elements into a scrollable container
	return container.NewVScroll(container.NewVBox(
		widget.NewLabel("Sampling Settings:"),
		samplingGrid,
		widget.NewSeparator(),
		streamingCheck,
		widget.NewLabel("Grammar-based Sampling:"),
		grammarEntry,
		widget.NewLabel("Logit Bias:"),
		logitBiasEntry,
		widget.NewLabel("Negative Prompt:"),
		negativePromptEntry,
		jsonModeCheck,
		speculativeNgramCheck,
		saveButton,
	))
}

func (t *TabLoad) saveAdvancedSettings() {
	// Collect all the settings
	settings := map[string]interface{}{
		"temperature":        t.temperatureSlider.Value,
		"top_k":              t.topKEntry.Text,
		"top_p":              t.topPSlider.Value,
		"min_p":              t.minPSlider.Value,
		"top_a":              t.topASlider.Value,
		"tfs":                t.tfsSlider.Value,
		"typical_p":          t.typicalPSlider.Value,
		"repetition_penalty": t.repetitionPenaltySlider.Value,
		"presence_penalty":   t.presencePenaltySlider.Value,
		"frequency_penalty":  t.frequencyPenaltySlider.Value,
		"mirostat_mode":      t.mirostatModeSelect.Selected,
		"mirostat_tau":       t.mirostatTauSlider.Value,
		"mirostat_eta":       t.mirostatEtaSlider.Value,
		"stream":             t.streamingCheck.Checked,
		"grammar_string":     t.grammarEntry.Text,
		"logit_bias":         t.logitBiasEntry.Text,
		"negative_prompt":    t.negativePromptEntry.Text,
		"json_mode":          t.jsonModeCheck.Checked,
		"speculative_ngram":  t.speculativeNgramCheck.Checked,
	}

	// Convert the settings to JSON
	jsonSettings, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to serialize settings: %v", err), t.window)
		return
	}

	// Save the settings to a file
	err = os.WriteFile("advanced_settings.json", jsonSettings, 0644)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to save settings: %v", err), t.window)
		return
	}

	dialog.ShowInformation("Settings Saved", "Advanced settings have been saved successfully.", t.window)
}

func (t *TabLoad) showSettingsDialog() {
	content := container.NewVBox(
		widget.NewCheck("Auto-connect on startup", func(checked bool) {
			t.saveAutoConnectSetting(checked)
		}),
	)

	dialog.ShowCustom("Settings", "Close", content, t.window)
}
