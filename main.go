package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Preset struct {
	Name               string
	MaxSeqLen          int
	OverrideBaseSeqLen int
	CacheSize          int
	GPUSplitAuto       bool
	GPUSplit           string
	RopeScale          float64
	RopeAlpha          float64
	CacheMode          string
	PromptTemplate     string
	NumExpertsPerToken int
	DraftModelName     string
	DraftRopeScale     float64
	DraftRopeAlpha     float64
	DraftCacheMode     string
	Fasttensors        bool
	AutosplitReserve   string
	ChunkSize          int
}

var (
	apiURL       string
	adminKey     string
	models       []string
	draftModels  []string
	loras        []string
	templates    []string
	overrides    []string
	presets      []Preset
	currentModel string
	currentLoras string

	modelsDropdown    *widget.Select
	lorasDropdown     *widget.Select
	templateDropdown  *widget.Select
	overrideDropdown  *widget.Select
	presetDropdown    *widget.Select
	currentModelLabel *widget.Label
	currentLorasLabel *widget.Label

	maxSeqLenEntry          *widget.Entry
	overrideBaseSeqLenEntry *widget.Entry
	cacheSizeEntry          *widget.Entry
	gpuSplitAutoCheck       *widget.Check
	gpuSplitEntry           *widget.Entry
	ropeScaleEntry          *widget.Entry
	ropeAlphaEntry          *widget.Entry
	cacheModeDropdown       *widget.Select
	promptTemplateEntry     *widget.Entry
	numExpertsPerTokenEntry *widget.Entry
	draftModelNameEntry     *widget.Entry
	draftRopeScaleEntry     *widget.Entry
	draftRopeAlphaEntry     *widget.Entry
	draftCacheModeDropdown  *widget.Select
	fasttensorsCheck        *widget.Check
	autosplitReserveEntry   *widget.Entry
	chunkSizeEntry          *widget.Entry
	repoIDEntry             *widget.Entry
	revisionEntry           *widget.Entry
	repoTypeDropdown        *widget.Select
	folderNameEntry         *widget.Entry
	includeEntry            *widget.Entry
	excludeEntry            *widget.Entry
	tokenEntry              *widget.Entry

	mu sync.Mutex
)

func main() {
	a := app.New()
	w := a.NewWindow("TabLoad")

	tabContainer := container.NewAppTabs()

	// Connection tab
	connectionTab := container.NewVBox()

	apiURLEntry := widget.NewEntry()
	apiURLEntry.SetPlaceHolder("TabbyAPI Endpoint URL")
	connectionTab.Add(apiURLEntry)

	adminKeyEntry := widget.NewPasswordEntry()
	adminKeyEntry.SetPlaceHolder("Admin Key")
	connectionTab.Add(adminKeyEntry)

	connectBtn := widget.NewButton("Connect", func() {
		url := apiURLEntry.Text
		key := adminKeyEntry.Text
		go connect(url, key, w)
	})
	connectionTab.Add(connectBtn)

	tabContainer.Append(container.NewTabItem("Connection", connectionTab))

	// Model tab
	modelTab := createModelTab()
	tabContainer.Append(container.NewTabItem("Model", modelTab))

	// Loras tab
	lorasTab := createLorasTab()
	tabContainer.Append(container.NewTabItem("Loras", lorasTab))

	// HF Downloader tab
	hfDownloaderTab := createHFDownloaderTab()
	tabContainer.Append(container.NewTabItem("HF Downloader", hfDownloaderTab))

	// Presets tab
	presetsTab := createPresetsTab()
	tabContainer.Append(container.NewTabItem("Presets", presetsTab))

	w.SetContent(tabContainer)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func createModelTab() *fyne.Container {
	modelTab := container.NewVBox()

	modelsDropdown = widget.NewSelect(models, func(selected string) {})
	modelTab.Add(modelsDropdown)

	maxSeqLenEntry = widget.NewEntry()
	maxSeqLenEntry.SetPlaceHolder("Max Sequence Length")
	modelTab.Add(maxSeqLenEntry)

	overrideBaseSeqLenEntry = widget.NewEntry()
	overrideBaseSeqLenEntry.SetPlaceHolder("Override Base Sequence Length")
	modelTab.Add(overrideBaseSeqLenEntry)

	cacheSizeEntry = widget.NewEntry()
	cacheSizeEntry.SetPlaceHolder("Cache Size")
	modelTab.Add(cacheSizeEntry)

	gpuSplitAutoCheck = widget.NewCheck("GPU Split Auto", nil)
	modelTab.Add(gpuSplitAutoCheck)

	gpuSplitEntry = widget.NewEntry()
	gpuSplitEntry.SetPlaceHolder("GPU Split")
	modelTab.Add(gpuSplitEntry)

	ropeScaleEntry = widget.NewEntry()
	ropeScaleEntry.SetPlaceHolder("Rope Scale")
	modelTab.Add(ropeScaleEntry)

	ropeAlphaEntry = widget.NewEntry()
	ropeAlphaEntry.SetPlaceHolder("Rope Alpha")
	modelTab.Add(ropeAlphaEntry)

	cacheModeDropdown = widget.NewSelect([]string{"Q4", "Q6", "Q8", "FP16"}, func(selected string) {})
	modelTab.Add(cacheModeDropdown)

	promptTemplateEntry = widget.NewEntry()
	promptTemplateEntry.SetPlaceHolder("Prompt Template")
	modelTab.Add(promptTemplateEntry)

	numExpertsPerTokenEntry = widget.NewEntry()
	numExpertsPerTokenEntry.SetPlaceHolder("Number of Experts per Token")
	modelTab.Add(numExpertsPerTokenEntry)

	draftModelNameEntry = widget.NewEntry()
	draftModelNameEntry.SetPlaceHolder("Draft Model Name")
	modelTab.Add(draftModelNameEntry)

	draftRopeScaleEntry = widget.NewEntry()
	draftRopeScaleEntry.SetPlaceHolder("Draft Rope Scale")
	modelTab.Add(draftRopeScaleEntry)

	draftRopeAlphaEntry = widget.NewEntry()
	draftRopeAlphaEntry.SetPlaceHolder("Draft Rope Alpha")
	modelTab.Add(draftRopeAlphaEntry)

	draftCacheModeDropdown = widget.NewSelect([]string{"Q4", "Q6", "Q8", "FP16"}, func(selected string) {})
	modelTab.Add(draftCacheModeDropdown)

	fasttensorsCheck = widget.NewCheck("Use Fasttensors", nil)
	modelTab.Add(fasttensorsCheck)

	autosplitReserveEntry = widget.NewEntry()
	autosplitReserveEntry.SetPlaceHolder("Autosplit Reserve")
	modelTab.Add(autosplitReserveEntry)

	chunkSizeEntry = widget.NewEntry()
	chunkSizeEntry.SetPlaceHolder("Chunk Size")
	modelTab.Add(chunkSizeEntry)

	loadModelBtn := widget.NewButton("Load Model", func() {
		loadModel()
	})
	modelTab.Add(loadModelBtn)

	unloadModelBtn := widget.NewButton("Unload Model", func() {
		unloadModel()
	})
	modelTab.Add(unloadModelBtn)

	currentModelLabel = widget.NewLabel("")
	modelTab.Add(currentModelLabel)

	return modelTab
}

func createLorasTab() *fyne.Container {
	lorasTab := container.NewVBox()

	lorasDropdown = widget.NewSelect(loras, func(selected string) {})
	lorasTab.Add(lorasDropdown)

	loadLorasBtn := widget.NewButton("Load Loras", func() {
		loadLoras([]string{lorasDropdown.Selected}, []float64{1.0})
	})
	lorasTab.Add(loadLorasBtn)

	unloadLorasBtn := widget.NewButton("Unload Loras", func() {
		unloadLoras()
	})
	lorasTab.Add(unloadLorasBtn)

	currentLorasLabel = widget.NewLabel("")
	lorasTab.Add(currentLorasLabel)

	return lorasTab
}

func createHFDownloaderTab() *fyne.Container {
	hfDownloaderTab := container.NewVBox()

	repoIDEntry = widget.NewEntry()
	repoIDEntry.SetPlaceHolder("Repo ID")
	hfDownloaderTab.Add(repoIDEntry)

	revisionEntry = widget.NewEntry()
	revisionEntry.SetPlaceHolder("Revision/Branch")
	hfDownloaderTab.Add(revisionEntry)

	repoTypeDropdown = widget.NewSelect([]string{"Model", "Lora"}, func(selected string) {})
	hfDownloaderTab.Add(repoTypeDropdown)

	folderNameEntry = widget.NewEntry()
	folderNameEntry.SetPlaceHolder("Folder Name")
	hfDownloaderTab.Add(folderNameEntry)

	includeEntry = widget.NewEntry()
	includeEntry.SetPlaceHolder("Include Patterns")
	hfDownloaderTab.Add(includeEntry)

	excludeEntry = widget.NewEntry()
	excludeEntry.SetPlaceHolder("Exclude Patterns")
	hfDownloaderTab.Add(excludeEntry)

	tokenEntry = widget.NewPasswordEntry()
	tokenEntry.SetPlaceHolder("HF Access Token")
	hfDownloaderTab.Add(tokenEntry)

	downloadBtn := widget.NewButton("Download", func() {
		download()
	})
	hfDownloaderTab.Add(downloadBtn)

	cancelDownloadBtn := widget.NewButton("Cancel", func() {
		cancelDownload()
	})
	hfDownloaderTab.Add(cancelDownloadBtn)

	return hfDownloaderTab
}

func createPresetsTab() *fyne.Container {
	presetsTab := container.NewVBox()

	presetDropdown = widget.NewSelect(getPresetNames(presets), func(selected string) {
		loadPreset(selected)
	})
	presetsTab.Add(presetDropdown)

	savePresetBtn := widget.NewButton("Save Preset", func() {
		savePreset()
	})
	presetsTab.Add(savePresetBtn)

	return presetsTab
}

func connect(url, key string, w fyne.Window) {
	mu.Lock()
	defer mu.Unlock()

	apiURL = url
	adminKey = key

	var err error
	models, err = fetchModels()
	if err != nil {
		showError(w, "Error fetching models", err)
		return
	}

	draftModels, err = fetchDraftModels()
	if err != nil {
		showError(w, "Error fetching draft models", err)
		return
	}

	loras, err = fetchLoras()
	if err != nil {
		showError(w, "Error fetching loras", err)
		return
	}

	templates, err = fetchTemplates()
	if err != nil {
		showError(w, "Error fetching templates", err)
		return
	}

	overrides, err = fetchOverrides()
	if err != nil {
		showError(w, "Error fetching overrides", err)
		return
	}

	presets = loadPresets()

	updateUI()

	currentModel, err = fetchCurrentModel()
	if err != nil {
		showError(w, "Error fetching current model", err)
		return
	}

	currentLoras, err = fetchCurrentLoras()
	if err != nil {
		showError(w, "Error fetching current loras", err)
		return
	}

	updateCurrentLabels()
}

func updateUI() {
	modelsDropdown.Options = models
	lorasDropdown.Options = loras
	templateDropdown.Options = templates
	overrideDropdown.Options = overrides
	presetDropdown.Options = getPresetNames(presets)
}

func updateCurrentLabels() {
	currentModelLabel.SetText(fmt.Sprintf("Current Model: %s", currentModel))
	currentLorasLabel.SetText(fmt.Sprintf("Current Loras: %s", currentLoras))
}

func showError(w fyne.Window, title string, err error) {
	dialog.ShowError(fmt.Errorf("%s: %v", title, err), w)
}

func fetchModels() ([]string, error) {
	url := apiURL + "/v1/model/list"
	body, err := makeHTTPRequest(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching models: %w", err)
	}

	var response struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling models JSON: %w", err)
	}

	models := make([]string, len(response.Data))
	for i, model := range response.Data {
		models[i] = model.ID
	}

	return models, nil
}

func fetchDraftModels() []string {
	// Fetch draft models from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/model/draft/list"
	body, err := makeHTTPRequest(url)
	if err != nil {
		fmt.Println("Error fetching draft models:", err)
		return nil
	}

	var response struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling draft models JSON:", err)
		return nil
	}

	draftModels := make([]string, len(response.Data))
	for i, model := range response.Data {
		draftModels[i] = model.ID
	}

	return draftModels
}

func fetchLoras() []string {
	// Fetch loras from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/lora/list"
	body, err := makeHTTPRequest(url)
	if err != nil {
		fmt.Println("Error fetching loras:", err)
		return nil
	}

	var response struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling loras JSON:", err)
		return nil
	}

	loras := make([]string, len(response.Data))
	for i, lora := range response.Data {
		loras[i] = lora.ID
	}

	return loras
}

func fetchTemplates() []string {
	// Fetch templates from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/template/list"
	body, err := makeHTTPRequest(url)
	if err != nil {
		fmt.Println("Error fetching templates:", err)
		return nil
	}

	var response struct {
		Data []string `json:"data"`
	}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling templates JSON:", err)
		return nil
	}

	return response.Data
}

func fetchOverrides() []string {
	// Fetch overrides from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/sampling/override/list"
	body, err := makeHTTPRequest(url)
	if err != nil {
		fmt.Println("Error fetching overrides:", err)
		return nil
	}

	var response struct {
		Presets []string `json:"presets"`
	}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling overrides JSON:", err)
		return nil
	}

	return response.Presets
}

func fetchCurrentModel() string {
	// Fetch current model from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/model"
	body, err := makeHTTPRequest(url)
	if err != nil {
		fmt.Println("Error fetching current model:", err)
		return ""
	}

	var response struct {
		ID         string `json:"id"`
		Parameters struct {
			MaxSeqLen int     `json:"max_seq_len"`
			CacheSize int     `json:"cache_size"`
			RopeScale float64 `json:"rope_scale"`
			RopeAlpha float64 `json:"rope_alpha"`
			Draft     struct {
				ID         string `json:"id"`
				Parameters struct {
					RopeScale float64 `json:"rope_scale"`
					RopeAlpha float64 `json:"rope_alpha"`
				} `json:"parameters"`
			} `json:"draft"`
		} `json:"parameters"`
	}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling current model JSON:", err)
		return ""
	}

	model := fmt.Sprintf("%s (context: %d, cache size: %d, rope scale: %.2f, rope alpha: %.2f)",
		response.ID, response.Parameters.MaxSeqLen, response.Parameters.CacheSize,
		response.Parameters.RopeScale, response.Parameters.RopeAlpha)

	if response.Parameters.Draft.ID != "" {
		model += fmt.Sprintf(" | %s (rope scale: %.2f, rope alpha: %.2f)",
			response.Parameters.Draft.ID, response.Parameters.Draft.Parameters.RopeScale,
			response.Parameters.Draft.Parameters.RopeAlpha)
	}

	return model
}

func fetchCurrentLoras() string {
	// Fetch current loras from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/lora"
	body, err := makeHTTPRequest(url)
	if err != nil {
		fmt.Println("Error fetching current loras:", err)
		return ""
	}

	var response struct {
		Data []struct {
			ID      string  `json:"id"`
			Scaling float64 `json:"scaling"`
		} `json:"data"`
	}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling current loras JSON:", err)
		return ""
	}

	loras := make([]string, len(response.Data))
	for i, lora := range response.Data {
		loras[i] = fmt.Sprintf("%s (scaling: %.2f)", lora.ID, lora.Scaling)
	}

	return strings.Join(loras, ", ")
}

func loadPresets() []Preset {
	presets := []Preset{}

	dir := filepath.Join(os.Getenv("HOME"), ".config", "tabload", "presets")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Error creating presets directory:", err)
		return presets
	}

	filePath := filepath.Join(dir, "presets.json")
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, create it with an empty array
			file, err = os.Create(filePath)
			if err != nil {
				fmt.Println("Error creating presets file:", err)
				return presets
			}
			file.WriteString("[]")
			file.Close()
			return presets
		}
		fmt.Println("Error opening presets file:", err)
		return presets
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&presets)
	if err != nil {
		fmt.Println("Error decoding presets JSON:", err)
	}

	return presets
}

func loadPreset(selectedPreset string) {
	for _, preset := range presets {
		if preset.Name == selectedPreset {
			// Load preset values into the UI
			modelsDropdown.SetSelected(preset.Name)
			maxSeqLenEntry.SetText(strconv.Itoa(preset.MaxSeqLen))
			overrideBaseSeqLenEntry.SetText(strconv.Itoa(preset.OverrideBaseSeqLen))
			cacheSizeEntry.SetText(strconv.Itoa(preset.CacheSize))
			gpuSplitAutoCheck.SetChecked(preset.GPUSplitAuto)
			gpuSplitEntry.SetText(preset.GPUSplit)
			ropeScaleEntry.SetText(fmt.Sprintf("%f", preset.RopeScale))
			ropeAlphaEntry.SetText(fmt.Sprintf("%f", preset.RopeAlpha))
			cacheModeDropdown.SetSelected(preset.CacheMode)
			promptTemplateEntry.SetText(preset.PromptTemplate)
			numExpertsPerTokenEntry.SetText(strconv.Itoa(preset.NumExpertsPerToken))
			draftModelNameEntry.SetText(preset.DraftModelName)
			draftRopeScaleEntry.SetText(fmt.Sprintf("%f", preset.DraftRopeScale))
			draftRopeAlphaEntry.SetText(fmt.Sprintf("%f", preset.DraftRopeAlpha))
			draftCacheModeDropdown.SetSelected(preset.DraftCacheMode)
			fasttensorsCheck.SetChecked(preset.Fasttensors)
			autosplitReserveEntry.SetText(preset.AutosplitReserve)
			chunkSizeEntry.SetText(strconv.Itoa(preset.ChunkSize))
			break
		}
	}
}

func savePreset() {
	// Save preset from the UI
	preset := Preset{
		Name:               "New Preset",
		MaxSeqLen:          1024,
		OverrideBaseSeqLen: 2048,
		CacheSize:          512,
		GPUSplitAuto:       true,
		GPUSplit:           "0.5,0.5",
		RopeScale:          1.0,
		RopeAlpha:          1.0,
		CacheMode:          "FP16",
		PromptTemplate:     "default",
		NumExpertsPerToken: 1,
		DraftModelName:     "",
		DraftRopeScale:     1.0,
		DraftRopeAlpha:     1.0,
		DraftCacheMode:     "FP16",
		Fasttensors:        false,
		AutosplitReserve:   "1024",
		ChunkSize:          128,
	}

	// Get values from model tab widgets
	// Get values from loras tab widgets
	// Get values from HF Downloader tab widgets

	presets := loadPresets()
	presets = append(presets, preset)

	file, err := os.Create("presets.json")
	if err != nil {
		fmt.Println("Error creating presets file:", err)
		return
	}
	defer file.Close()

	jsonEncoder := json.NewEncoder(file)
	err = jsonEncoder.Encode(presets)
	if err != nil {
		fmt.Println("Error encoding presets JSON:", err)
		return
	}

	fmt.Println("Saved preset:", preset.Name)
}

func makeHTTPRequest(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-api-key", adminKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func loadModel(modelName string, maxSeqLen int, overrideBaseSeqLen int, cacheSize int, gpuSplitAuto bool, gpuSplit string, modelRopeScale float64, modelRopeAlpha float64, cacheMode string, promptTemplate string, numExpertsPerToken int, draftModelName string, draftRopeScale float64, draftRopeAlpha float64, draftCacheMode string, fasttensors bool, autosplitReserve string, chunkSize int) {
	// Load model from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/model/load"

	request := map[string]interface{}{
		"name":                  modelName,
		"max_seq_len":           maxSeqLen,
		"override_base_seq_len": overrideBaseSeqLen,
		"cache_size":            cacheSize,
		"gpu_split_auto":        gpuSplitAuto,
		"gpu_split":             gpuSplit,
		"rope_scale":            modelRopeScale,
		"rope_alpha":            modelRopeAlpha,
		"cache_mode":            cacheMode,
		"prompt_template":       promptTemplate,
		"num_experts_per_token": numExpertsPerToken,
		"fasttensors":           fasttensors,
		"autosplit_reserve":     autosplitReserve,
		"chunk_size":            chunkSize,
	}

	if draftModelName != "" {
		draftRequest := map[string]interface{}{
			"draft_model_name": draftModelName,
			"draft_rope_scale": draftRopeScale,
			"draft_rope_alpha": draftRopeAlpha,
			"draft_cache_mode": draftCacheMode,
		}
		request["draft"] = draftRequest
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling request JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Model loaded successfully:", string(body))
}

func loadLoras(loras []string, scalings []float64) {
	// Load loras from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/lora/load"

	loadList := make([]map[string]interface{}, len(loras))
	for i, lora := range loras {
		loadList[i] = map[string]interface{}{
			"name":    lora,
			"scaling": scalings[i],
		}
	}

	request := map[string]interface{}{
		"loras": loadList,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling request JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Loras loaded successfully:", string(body))
}

func unloadModel() {
	// Unload model from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/model/unload"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Model unloaded successfully")
}

func unloadLoras() {
	// Unload loras from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/lora/unload"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Loras unloaded successfully")
}

func loadTemplate(promptTemplate string) {
	// Load template from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/template/switch"

	request := map[string]string{
		"name": promptTemplate,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling request JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Prompt template loaded successfully")
}

func unloadTemplate() {
	// Unload template from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/template/unload"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Prompt template unloaded successfully")
}

func loadOverride(samplerOverride string) {
	// Load override from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/sampling/override/switch"

	request := map[string]string{
		"preset": samplerOverride,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling request JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Sampler override loaded successfully")
}

func unloadOverride() {
	// Unload override from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/sampling/override/unload"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Sampler override unloaded successfully")
}

func download(repoID string, revision string, repoType string, folderName string, token string, include string, exclude string) {
	// Download from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/download"

	includePatterns := strings.Split(include, ",")
	excludePatterns := strings.Split(exclude, ",")

	request := map[string]interface{}{
		"repo_id":     repoID,
		"revision":    revision,
		"repo_type":   repoType,
		"folder_name": folderName,
		"token":       token,
		"include":     includePatterns,
		"exclude":     excludePatterns,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling request JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var response struct {
		DownloadPath string `json:"download_path"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response JSON:", err)
		return
	}

	fmt.Println("Download completed successfully. Path:", response.DownloadPath)
}

func cancelDownload() {
	// Cancel download from API
	// Replace the URL with the actual API endpoint
	url := apiURL + "/v1/download/cancel"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("X-admin-key", adminKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Download canceled successfully")
}

func parseInt(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func parseFloat64(value string) (float64, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return floatValue, nil
}

func parseBool(value string) (bool, error) {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return boolValue, nil
}

func getPresetNames(presets []Preset) []string {
	names := make([]string, len(presets))
	for i, preset := range presets {
		names[i] = preset.Name
	}
	return names
}
