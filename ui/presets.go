package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sammcj/tabload/logging"
	"github.com/sammcj/tabload/utils"
)

var cachedPresets []Preset

func (t *TabLoad) loadPresetsFromStorage() ([]Preset, error) {
	if cachedPresets != nil {
		return cachedPresets, nil
	}

	defaultPresets := []Preset{
		{
			Name:           "Default Preset",
			CacheMode:      "Q4",
			DraftCacheMode: "Q4",
		},
	}

	logging.Info(fmt.Sprintf("Loading presets from: %s", presetsFilePath))

	data, err := os.ReadFile(presetsFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logging.Info("Presets file not found. Using default presets.")
			cachedPresets = defaultPresets
			return cachedPresets, nil
		}
		return nil, fmt.Errorf("error reading presets file: %w", err)
	}

	if len(data) == 0 {
		logging.Info("Presets file is empty. Using default presets.")
		cachedPresets = defaultPresets
		return cachedPresets, nil
	}

	var presetsWrapper struct {
		Presets []Preset `json:"presets"`
	}

	err = json.Unmarshal(data, &presetsWrapper)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling presets: %w", err)
	}

	loadedPresets := presetsWrapper.Presets

	// Merge loaded presets with default presets
	presetMap := make(map[string]Preset)
	for _, p := range defaultPresets {
		presetMap[p.Name] = p
	}
	for _, p := range loadedPresets {
		presetMap[p.Name] = p
	}

	mergedPresets := make([]Preset, 0, len(presetMap))
	for _, p := range presetMap {
		mergedPresets = append(mergedPresets, p)
	}

	cachedPresets = mergedPresets
	logging.Info(fmt.Sprintf("Loaded %d presets", len(mergedPresets)))

	return cachedPresets, nil
}

func (t *TabLoad) savePresetsToStorage(presets []Preset) error {
	presetsWrapper := struct {
		Presets []Preset `json:"presets"`
	}{
		Presets: presets,
	}

	data, err := json.MarshalIndent(presetsWrapper, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling presets: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(presetsFilePath), 0755); err != nil {
		return fmt.Errorf("error creating presets directory: %w", err)
	}

	if err := os.WriteFile(presetsFilePath, data, 0644); err != nil {
		return fmt.Errorf("error writing presets file: %w", err)
	}

	return nil
}

func (t *TabLoad) validatePreset(preset Preset) error {
	if preset.Name == "" {
		return fmt.Errorf("preset name cannot be empty")
	}
	// Add more validation rules as needed
	return nil
}

func (t *TabLoad) savePresetToStorage(preset Preset) error {
	if err := t.validatePreset(preset); err != nil {
		return fmt.Errorf("invalid preset: %w", err)
	}

	presets, err := t.loadPresetsFromStorage()
	if err != nil {
		return fmt.Errorf("error loading existing presets: %w", err)
	}

	updated := false
	for i, p := range presets {
		if p.Name == preset.Name {
			presets[i] = preset
			updated = true
			break
		}
	}
	if !updated {
		presets = append(presets, preset)
	}

	return t.savePresetsToStorage(presets)
}

func (t *TabLoad) loadPresetFromStorage(presetName string) (*Preset, error) {
	presets, err := t.loadPresetsFromStorage()
	if err != nil {
		return nil, fmt.Errorf("error loading presets: %w", err)
	}

	for _, preset := range presets {
		if preset.Name == presetName {
			return &preset, nil
		}
	}

	return nil, fmt.Errorf("preset not found: %s", presetName)
}

func (t *TabLoad) deletePresetFromStorage(presetName string) error {
	presets, err := t.loadPresetsFromStorage()
	if err != nil {
		return fmt.Errorf("error loading existing presets: %w", err)
	}

	newPresets := make([]Preset, 0, len(presets))
	found := false
	for _, p := range presets {
		if p.Name != presetName {
			newPresets = append(newPresets, p)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("preset not found: %s", presetName)
	}

	// take a backup of the presets file before deleting
	if err := os.Rename(presetsFilePath, presetsFilePath+".bak"); err != nil {
		return fmt.Errorf("error taking backup of presets file: %w", err)
	}

	return t.savePresetsToStorage(newPresets)
}

// func (t *TabLoad) refreshPresetList() {
// 	if t.presetDropdown != nil && len(t.presetDropdown.Options) > 1 {
// 		// Presets are already loaded, no need to refresh
// 		return
// 	}

// 	presets, err := t.loadPresetsFromStorage()
// 	if err != nil {
// 		logging.Error("Error loading presets", err)
// 		dialog.ShowError(err, t.window)
// 		return
// 	}

// 	presetNames := make([]string, 0, len(presets)+1)
// 	presetNames = append(presetNames, "(Select one)") // Add default option
// 	for _, preset := range presets {
// 		presetNames = append(presetNames, preset.Name)
// 	}

// 	t.presetDropdown.Options = presetNames
// 	t.presetDropdown.Refresh()

// 	if len(presetNames) > 1 {
// 		t.presetDropdown.SetSelected(presetNames[1]) // Select the first actual preset
// 	} else {
// 		t.presetDropdown.SetSelected(presetNames[0]) // Select "(Select one)" if no presets
// 	}

// 	logging.Info(fmt.Sprintf("Refreshed preset list with %d presets", len(presetNames)-1))
// }

func (t *TabLoad) refreshPresetList() {
	presets, err := t.loadPresetsFromStorage()
	if err != nil {
		logging.Error("Error loading presets", err)
		dialog.ShowError(err, t.window)
		return
	}

	presetNames := make([]string, 0, len(presets)+1)
	presetNames = append(presetNames, "(Select one)")
	for _, preset := range presets {
		presetNames = append(presetNames, preset.Name)
	}

	t.presetDropdown.Options = presetNames
	t.presetDropdown.Refresh()

	if len(presetNames) > 1 {
		t.presetDropdown.SetSelected(presetNames[1])
	} else {
		t.presetDropdown.SetSelected(presetNames[0])
	}

	logging.Info(fmt.Sprintf("Refreshed preset list with %d presets", len(presetNames)-1))
}

func (t *TabLoad) handleSavePreset() {
	presetName := widget.NewEntry()
	presetName.SetPlaceHolder("Enter preset name")

	dialog.ShowForm("Save Preset", "Save", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Preset Name", presetName),
	}, func(save bool) {
		if !save || presetName.Text == "" {
			return
		}

		preset := t.createPresetFromFields()
		preset.Name = presetName.Text

		if err := t.savePresetToStorage(preset); err != nil {
			logging.Error("Failed to save preset", err)
			dialog.ShowError(fmt.Errorf("failed to save preset: %w", err), t.window)
			return
		}

		t.refreshPresetList()
		logging.Info(fmt.Sprintf("Preset '%s' saved successfully", preset.Name))
		dialog.ShowInformation("Success", "Preset saved successfully", t.window)
	}, t.window)
}

func (t *TabLoad) applyPresetToFields(preset *Preset) {
	// Helper function to set entry text and enable if value is present
	setEntryText := func(entry *widget.Entry, checkbox *widget.Check, value interface{}) {
		if entry == nil || checkbox == nil {
			return
		}
		if value != nil && value != "" {
			entry.SetText(fmt.Sprintf("%v", value))
			checkbox.SetChecked(true)
			entry.Enable()
		} else {
			entry.SetText("")
			checkbox.SetChecked(false)
			entry.Disable()
		}
	}

	// Apply preset values to fields
	if t.modelsDropdown != nil {
		t.modelsDropdown.SetSelected(preset.Name)
	}
	setEntryText(t.maxSeqLenEntry, t.maxSeqLenCheck, preset.MaxSeqLen)
	setEntryText(t.overrideBaseSeqLenEntry, t.overrideBaseSeqLenCheck, preset.OverrideBaseSeqLen)
	setEntryText(t.cacheSizeEntry, t.cacheSizeCheck, preset.CacheSize)
	if t.gpuSplitAutoCheck != nil {
		t.gpuSplitAutoCheck.SetChecked(preset.GPUSplitAuto)
	}
	setEntryText(t.gpuSplitEntry, t.gpuSplitCheck, preset.GPUSplit)
	setEntryText(t.ropeScaleEntry, t.ropeScaleCheck, preset.RopeScale)
	setEntryText(t.ropeAlphaEntry, t.ropeAlphaCheck, preset.RopeAlpha)
	if t.cacheModeDropdown != nil {
		t.cacheModeDropdown.SetSelected(preset.CacheMode)
	}
	setEntryText(t.promptTemplateEntry, t.promptTemplateCheck, preset.PromptTemplate)
	setEntryText(t.numExpertsPerTokenEntry, t.numExpertsPerTokenCheck, preset.NumExpertsPerToken)
	setEntryText(t.draftModelNameEntry, t.draftModelNameCheck, preset.DraftModelName)
	setEntryText(t.draftRopeScaleEntry, t.draftRopeScaleCheck, preset.DraftRopeScale)
	setEntryText(t.draftRopeAlphaEntry, t.draftRopeAlphaCheck, preset.DraftRopeAlpha)
	if t.draftCacheModeDropdown != nil {
		t.draftCacheModeDropdown.SetSelected(preset.DraftCacheMode)
	}
	if t.fasttensorsCheck != nil {
		t.fasttensorsCheck.SetChecked(preset.Fasttensors)
	}
	setEntryText(t.autosplitReserveEntry, t.autosplitReserveCheck, preset.AutosplitReserve)
	setEntryText(t.chunkSizeEntry, t.chunkSizeCheck, preset.ChunkSize)
}

func (t *TabLoad) createPresetFromFields() Preset {
	preset := Preset{
		Name:             t.modelsDropdown.Selected,
		GPUSplitAuto:     t.gpuSplitAutoCheck.Checked,
		GPUSplit:         t.gpuSplitEntry.Text,
		CacheMode:        t.cacheModeDropdown.Selected,
		PromptTemplate:   utils.ParseStringPointer(t.promptTemplateEntry.Text),
		DraftCacheMode:   t.draftCacheModeDropdown.Selected,
		Fasttensors:      t.fasttensorsCheck.Checked,
		AutosplitReserve: t.autosplitReserveEntry.Text,
	}

	preset.MaxSeqLen = utils.ParseIntPointer(t.maxSeqLenEntry.Text)
	preset.OverrideBaseSeqLen = utils.ParseIntPointer(t.overrideBaseSeqLenEntry.Text)
	preset.CacheSize = utils.ParseIntPointer(t.cacheSizeEntry.Text)
	preset.RopeScale = utils.ParseFloat64Pointer(t.ropeScaleEntry.Text)
	preset.RopeAlpha = utils.ParseFloat64Pointer(t.ropeAlphaEntry.Text)
	preset.NumExpertsPerToken = utils.ParseIntPointer(t.numExpertsPerTokenEntry.Text)
	preset.DraftModelName = utils.ParseStringPointer(t.draftModelNameEntry.Text)
	preset.DraftRopeScale = utils.ParseFloat64Pointer(t.draftRopeScaleEntry.Text)
	preset.DraftRopeAlpha = utils.ParseFloat64Pointer(t.draftRopeAlphaEntry.Text)
	preset.ChunkSize = utils.ParseIntPointer(t.chunkSizeEntry.Text)

	return preset
}

func (t *TabLoad) handleLoadPreset(selectedPreset string) {
	if selectedPreset == "" || selectedPreset == "(Select one)" {
		t.clearAllFields()
		return
	}

	preset, err := t.loadPresetFromStorage(selectedPreset)
	if err != nil {
		logging.Error("Failed to load preset", err)
		dialog.ShowError(err, t.window)
		return
	}

	if preset.MaxSeqLen != nil {
		t.maxSeqLenEntry.SetText(fmt.Sprintf("%d", *preset.MaxSeqLen))
	} else {
		t.maxSeqLenEntry.SetText("")
	}

	t.applyPresetToFields(preset)
	logging.Info(fmt.Sprintf("Preset '%s' loaded successfully", selectedPreset))
}

func (t *TabLoad) clearAllFields() {
	// Clear and disable all fields
	entries := []*widget.Entry{
		t.maxSeqLenEntry, t.overrideBaseSeqLenEntry, t.cacheSizeEntry,
		t.gpuSplitEntry, t.ropeScaleEntry, t.ropeAlphaEntry,
		t.promptTemplateEntry, t.numExpertsPerTokenEntry,
		t.draftModelNameEntry, t.draftRopeScaleEntry, t.draftRopeAlphaEntry,
		t.autosplitReserveEntry, t.chunkSizeEntry,
	}
	checkboxes := []*widget.Check{
		t.maxSeqLenCheck, t.overrideBaseSeqLenCheck, t.cacheSizeCheck,
		t.gpuSplitCheck, t.ropeScaleCheck, t.ropeAlphaCheck,
		t.promptTemplateCheck, t.numExpertsPerTokenCheck,
		t.draftModelNameCheck, t.draftRopeScaleCheck, t.draftRopeAlphaCheck,
		t.autosplitReserveCheck, t.chunkSizeCheck,
	}
	for i, entry := range entries {
		entry.SetText("")
		entry.Disable()
		checkboxes[i].SetChecked(false)
	}
	t.gpuSplitAutoCheck.SetChecked(false)
	t.fasttensorsCheck.SetChecked(false)
	t.cacheModeDropdown.SetSelected("")
	t.draftCacheModeDropdown.SetSelected("")
}

func (t *TabLoad) handleDeletePreset() {
	if t.presetDropdown.Selected == "" {
		logging.Warn("No preset selected for deletion")
		dialog.ShowInformation("Error", "Please select a preset to delete", t.window)
		return
	}
	dialog.ShowConfirm("Delete Preset", "Are you sure you want to delete this preset?", func(confirm bool) {
		if !confirm {
			return
		}

		if err := t.deletePresetFromStorage(t.presetDropdown.Selected); err != nil {
			logging.Error("Failed to delete preset", err)
			dialog.ShowError(err, t.window)
			return
		}
		t.refreshPresetList()
		logging.Info(fmt.Sprintf("Preset '%s' deleted successfully", t.presetDropdown.Selected))
		dialog.ShowInformation("Success", "Preset deleted successfully", t.window)
	}, t.window)
}

func (t *TabLoad) buildPresetTab() fyne.CanvasObject {
	t.presetDropdown = widget.NewSelect([]string{}, t.handleLoadPreset)

	saveButton := widget.NewButton("Save Preset", t.handleSavePreset)
	deleteButton := widget.NewButton("Delete Preset", t.handleDeletePreset)

	return container.NewVBox(t.presetDropdown, saveButton, deleteButton)
}
