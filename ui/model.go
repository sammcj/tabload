package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sammcj/tabload/logging"
	"github.com/sammcj/tabload/utils"
)

func (t *TabLoad) initialiseUIElements() {
	// Ensure the client is initialised
	if t.client == nil {
		logging.Error("Client is not initialised", nil)
		dialog.ShowError(fmt.Errorf("client is not initialised"), t.window)
		return
	}

	// Initialise dropdown and entry widgets
	t.modelsDropdown = widget.NewSelect([]string{}, func(selected string) {})
	t.maxSeqLenEntry = widget.NewEntry()
	t.overrideBaseSeqLenEntry = widget.NewEntry()
	t.cacheSizeEntry = widget.NewEntry()
	t.gpuSplitEntry = widget.NewEntry()
	t.ropeScaleEntry = widget.NewEntry()
	t.ropeAlphaEntry = widget.NewEntry()
	t.promptTemplateEntry = widget.NewEntry()
	t.numExpertsPerTokenEntry = widget.NewEntry()
	t.draftModelNameEntry = widget.NewEntry()
	t.draftRopeScaleEntry = widget.NewEntry()
	t.draftRopeAlphaEntry = widget.NewEntry()
	t.autosplitReserveEntry = widget.NewEntry()
	t.chunkSizeEntry = widget.NewEntry()

	// Initialise checkbox widgets
	t.gpuSplitAutoCheck = widget.NewCheck("", func(bool) {})
	t.fasttensorsCheck = widget.NewCheck("", func(bool) {})
	t.maxSeqLenCheck = widget.NewCheck("", nil)
	t.overrideBaseSeqLenCheck = widget.NewCheck("", nil)
	t.cacheSizeCheck = widget.NewCheck("", nil)
	t.gpuSplitCheck = widget.NewCheck("", nil)
	t.ropeScaleCheck = widget.NewCheck("", nil)
	t.ropeAlphaCheck = widget.NewCheck("", nil)
	t.promptTemplateCheck = widget.NewCheck("", nil)
	t.numExpertsPerTokenCheck = widget.NewCheck("", nil)
	t.draftModelNameCheck = widget.NewCheck("", nil)
	t.draftRopeScaleCheck = widget.NewCheck("", nil)
	t.draftRopeAlphaCheck = widget.NewCheck("", nil)
	t.autosplitReserveCheck = widget.NewCheck("", nil)
	t.chunkSizeCheck = widget.NewCheck("", nil)

	// Initialise other widgets
	t.cacheModeDropdown = widget.NewSelect([]string{"Q4", "Q6", "Q8", "FP16"}, func(selected string) {})
	t.draftCacheModeDropdown = widget.NewSelect([]string{"Q4", "Q6", "Q8", "FP16"}, func(selected string) {})
	t.presetDropdown = widget.NewSelect([]string{}, t.handleLoadPreset)
	t.loadModelButton = widget.NewButton("Load Model", t.handleLoadModel)
	t.unloadModelButton = widget.NewButton("Unload Model", t.handleUnloadModel)
	t.currentModelLabel = widget.NewLabel("")
	t.savePresetButton = widget.NewButton("Save Preset", t.handleSavePreset)
	t.deletePresetButton = widget.NewButton("Delete Preset", t.handleDeletePreset)

	// Set placeholder text for entries
	t.maxSeqLenEntry.SetPlaceHolder("Enter max sequence length")
	t.overrideBaseSeqLenEntry.SetPlaceHolder("Enter override base seq length")
	t.cacheSizeEntry.SetPlaceHolder("Enter cache size")
	t.gpuSplitEntry.SetPlaceHolder("Enter GPU split")
	t.ropeScaleEntry.SetPlaceHolder("Enter rope scale")
	t.ropeAlphaEntry.SetPlaceHolder("Enter rope alpha")
	t.promptTemplateEntry.SetPlaceHolder("Enter prompt template")
	t.numExpertsPerTokenEntry.SetPlaceHolder("Enter num experts per token")
	t.draftModelNameEntry.SetPlaceHolder("Enter draft model name")
	t.draftRopeScaleEntry.SetPlaceHolder("Enter draft rope scale")
	t.draftRopeAlphaEntry.SetPlaceHolder("Enter draft rope alpha")
	t.autosplitReserveEntry.SetPlaceHolder("Enter autosplit reserve")
	t.chunkSizeEntry.SetPlaceHolder("Enter chunk size")

	t.refreshPresetList()

	// Fetch and set prompt templates
	// only if we have a client
	if t.client != nil {
		templates, err := t.client.FetchTemplates()
		if err != nil {
			logging.Error("Error fetching templates", err)
			dialog.ShowError(err, t.window)
			return
		}
		t.promptTemplateEntry.SetText(strings.Join(templates, ", "))
	}

	// Disable all entry fields and set checkboxes to unchecked
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
		entry.Disable()
		checkboxes[i].SetChecked(false)
		checkboxes[i].OnChanged = func(checked bool) {
			if checked {
				entry.Enable()
			} else {
				entry.Disable()
				entry.SetText("")
			}
		}
	}
}

// Update the handleLoadModel function to use the new UI elements
func (t *TabLoad) handleLoadModel() {
	params := make(map[string]interface{})

	// Only add parameters that have been set by the user
	if t.modelsDropdown.Selected != "" {
		params["name"] = t.modelsDropdown.Selected
	}

	if t.maxSeqLenCheck.Checked {
		if maxSeqLen, err := strconv.Atoi(t.maxSeqLenEntry.Text); err == nil {
			params["max_seq_len"] = maxSeqLen
		}
	}

	if t.overrideBaseSeqLenCheck.Checked {
		params["override_base_seq_len"] = utils.ParseIntOrZero(t.overrideBaseSeqLenEntry.Text)
	}

	if t.cacheSizeCheck.Checked {
		params["cache_size"] = utils.ParseIntOrZero(t.cacheSizeEntry.Text)
	}

	if t.gpuSplitAutoCheck.Checked {
		params["gpu_split_auto"] = true
	} else if t.gpuSplitCheck.Checked {
		params["gpu_split"] = t.gpuSplitEntry.Text
	}

	if t.ropeScaleCheck.Checked {
		params["rope_scale"] = utils.ParseFloatOrZero(t.ropeScaleEntry.Text)
	}

	if t.ropeAlphaCheck.Checked {
		params["rope_alpha"] = utils.ParseFloatOrZero(t.ropeAlphaEntry.Text)
	}

	if t.cacheModeDropdown.Selected != "" {
		params["cache_mode"] = t.cacheModeDropdown.Selected
	}

	if t.promptTemplateCheck.Checked && t.promptTemplateEntry.Text != "" {
		params["prompt_template"] = t.promptTemplateEntry.Text
	}

	if t.numExpertsPerTokenCheck.Checked {
		params["num_experts_per_token"] = utils.ParseIntOrZero(t.numExpertsPerTokenEntry.Text)
	}

	if t.fasttensorsCheck.Checked {
		params["fasttensors"] = true
	}

	if t.autosplitReserveCheck.Checked {
		params["autosplit_reserve"] = t.autosplitReserveEntry.Text
	}

	if t.chunkSizeCheck.Checked {
		params["chunk_size"] = utils.ParseIntOrZero(t.chunkSizeEntry.Text)
	}

	// Handle draft model parameters
	if t.draftModelNameCheck.Checked && t.draftModelNameEntry.Text != "" {
		draftParams := make(map[string]interface{})
		draftParams["draft_model_name"] = t.draftModelNameEntry.Text

		if t.draftRopeScaleCheck.Checked {
			draftParams["draft_rope_scale"] = utils.ParseFloatOrZero(t.draftRopeScaleEntry.Text)
		}

		if t.draftRopeAlphaCheck.Checked {
			draftParams["draft_rope_alpha"] = utils.ParseFloatOrZero(t.draftRopeAlphaEntry.Text)
		}

		if t.draftCacheModeDropdown.Selected != "" {
			draftParams["draft_cache_mode"] = t.draftCacheModeDropdown.Selected
		}

		params["draft"] = draftParams
	}

	err := t.client.LoadModel(t.modelsDropdown.Selected, params)
	if err != nil {
		logging.Error("Error loading model", err)
		dialog.ShowError(err, t.window)
		return
	}

	t.refreshCurrentModel()
	logging.Info("Model loaded successfully")
}

func (t *TabLoad) buildModelTab() fyne.CanvasObject {
	t.initialiseUIElements()

	// Create a new form
	t.form = &widget.Form{
		Items: []*widget.FormItem{},
	}

	// Fetch templates
	templates := t.loadTemplates()
	templateNames := make([]string, 0, len(templates)+1)
	for name := range templates {
		templateNames = append(templateNames, name)
	}
	// sort the template names alphabetically, but if they contain the word "server" they should be at the end
	sort.Slice(templateNames, func(i, j int) bool {
		if strings.Contains(templateNames[i], "server") && !strings.Contains(templateNames[j], "server") {
			return false
		}
		if !strings.Contains(templateNames[i], "server") && strings.Contains(templateNames[j], "server") {
			return true
		}
		return templateNames[i] < templateNames[j]
	})

	templateNames = append(templateNames, "Create New...")

	t.promptTemplateDropdown = widget.NewSelect(templateNames, func(selected string) {
		if selected == "Create New..." {
			t.showCreateTemplateDialog()
		} else {
			t.promptTemplateEntry.SetText(strings.TrimSuffix(selected, " (server)"))
		}
	})

	deleteButton := widget.NewButton("Delete", func() {
		if t.promptTemplateDropdown.Selected != "" && t.promptTemplateDropdown.Selected != "Create New..." {
			t.deleteTemplate(t.promptTemplateDropdown.Selected)
		}
	})

	templateContainer := container.NewBorder(nil, nil, nil, deleteButton, t.promptTemplateDropdown)

	if t.promptTemplateEntry.Text != "" {
		t.promptTemplateDropdown.SetSelected(t.promptTemplateEntry.Text)
	}

	// Add rows to the form
	t.addFormRow("Model", t.modelsDropdown)
	t.addFormRow("Max Sequence Length", t.createCheckboxEntry(t.maxSeqLenCheck, t.maxSeqLenEntry))
	t.addFormRow("Override Base Seq Length", t.createCheckboxEntry(t.overrideBaseSeqLenCheck, t.overrideBaseSeqLenEntry))
	t.addFormRow("Cache Size", t.createCheckboxEntry(t.cacheSizeCheck, t.cacheSizeEntry))
	t.addFormRow("GPU Split Auto", t.gpuSplitAutoCheck)
	t.addFormRow("GPU Split", t.createCheckboxEntry(t.gpuSplitCheck, t.gpuSplitEntry))
	t.addFormRow("Rope Scale", t.createCheckboxEntry(t.ropeScaleCheck, t.ropeScaleEntry))
	t.addFormRow("Rope Alpha", t.createCheckboxEntry(t.ropeAlphaCheck, t.ropeAlphaEntry))
	t.addFormRow("Cache Mode", t.cacheModeDropdown)
	t.addFormRow("Prompt Template", templateContainer)
	t.addFormRow("Num Experts Per Token", t.createCheckboxEntry(t.numExpertsPerTokenCheck, t.numExpertsPerTokenEntry))
	t.addFormRow("Draft Model Name", t.createCheckboxEntry(t.draftModelNameCheck, t.draftModelNameEntry))
	t.addFormRow("Draft Rope Scale", t.createCheckboxEntry(t.draftRopeScaleCheck, t.draftRopeScaleEntry))
	t.addFormRow("Draft Rope Alpha", t.createCheckboxEntry(t.draftRopeAlphaCheck, t.draftRopeAlphaEntry))
	t.addFormRow("Draft Cache Mode", t.draftCacheModeDropdown)
	t.addFormRow("Use Fasttensors", t.fasttensorsCheck)
	t.addFormRow("Autosplit Reserve", t.createCheckboxEntry(t.autosplitReserveCheck, t.autosplitReserveEntry))
	t.addFormRow("Chunk Size", t.createCheckboxEntry(t.chunkSizeCheck, t.chunkSizeEntry))

	// Create containers for presets and buttons
	presetContainer := container.NewHBox(t.presetDropdown, t.savePresetButton, t.deletePresetButton)
	buttonsContainer := container.NewHBox(t.loadModelButton, t.unloadModelButton)

	// Combine all elements
	return container.NewVBox(
		presetContainer,
		t.form,
		buttonsContainer,
		t.currentModelLabel,
	)
}

func (t *TabLoad) buildLorasTab() fyne.CanvasObject {
	t.lorasDropdown = widget.NewSelect([]string{}, func(selected string) {})
	t.loadLorasButton = widget.NewButton("Load LoRAs", t.handleLoadLoras)
	t.unloadLorasButton = widget.NewButton("Unload LoRAs", t.handleUnloadLoras)

	t.currentLorasLabel = widget.NewLabel("")

	return container.NewVBox(
		t.lorasDropdown,
		t.loadLorasButton,
		t.unloadLorasButton,
		t.currentLorasLabel,
	)
}

func (t *TabLoad) handleUnloadModel() {
	err := t.client.UnloadModel()
	if err != nil {
		// Handle error (e.g., show an error dialog)
		fmt.Println("Error unloading model:", err)
		return
	}

	t.refreshCurrentModel()
}

func (t *TabLoad) handleLoadLoras() {
	selectedLoras := []string{t.lorasDropdown.Selected}
	err := t.client.LoadLoras(selectedLoras, []float64{1.0})
	if err != nil {
		// Handle error (e.g., show an error dialog)
		fmt.Println("Error loading LoRAs:", err)
		return
	}

	t.refreshCurrentLoras()
}

func (t *TabLoad) handleUnloadLoras() {
	err := t.client.UnloadLoras()
	if err != nil {
		// Handle error (e.g., show an error dialog)
		fmt.Println("Error unloading loras:", err)
		return
	}

	t.refreshCurrentLoras()
}

func (t *TabLoad) refreshCurrentModel() {
	if t.client == nil {
		logging.Warn("Cannot refresh current model: client is nil")
		t.updateModelInfoContainer([][]string{{"No connection established", ""}})
		return
	}

	currentModel, err := t.client.FetchCurrentModel()
	if err != nil {
		logging.Error("Error fetching current model", err)
		t.updateModelInfoContainer([][]string{{"Error fetching current model", ""}})
		return
	}

	if currentModel == nil {
		logging.Warn("Received nil current model")
		t.updateModelInfoContainer([][]string{{"No model loaded", ""}})
		return
	}

	modelInfo := [][]string{
		{"Model", currentModel.ID},
		{"Max Sequence Length", fmt.Sprintf("%d", currentModel.Parameters.MaxSeqLen)},
		{"Cache Size", fmt.Sprintf("%d", currentModel.Parameters.CacheSize)},
		{"Rope Scale", fmt.Sprintf("%.2f", currentModel.Parameters.RopeScale)},
		{"Rope Alpha", fmt.Sprintf("%.2f", currentModel.Parameters.RopeAlpha)},
	}

	if currentModel.Parameters.Draft != nil {
		modelInfo = append(modelInfo,
			[]string{"Draft Model", currentModel.Parameters.Draft.ID},
			[]string{"Draft Rope Scale", fmt.Sprintf("%.2f", currentModel.Parameters.Draft.Parameters.RopeScale)},
			[]string{"Draft Rope Alpha", fmt.Sprintf("%.2f", currentModel.Parameters.Draft.Parameters.RopeAlpha)},
		)
	}

	t.updateModelInfoContainer(modelInfo)
}

func (t *TabLoad) addFormRow(label string, widget fyne.CanvasObject) {
	// t.form.Add(widget.NewLabel(label))
	// t.form.Add(widget)
	t.form.Append(label, widget)
	t.form.Refresh()
}

func (t *TabLoad) createCheckboxEntry(checkbox *widget.Check, entry *widget.Entry) *fyne.Container {
	return container.NewBorder(nil, nil, checkbox, nil, entry)
}

func (t *TabLoad) handleModelSelection(modelName string) {
	// Clear existing fields
	t.clearAllFields()

	// Load default parameters
	t.LoadDefaultParams()

	// If a model is selected, you might want to load model-specific parameters here
	if modelName != "" {
		// Load model-specific parameters if available
	}
}
