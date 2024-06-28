package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sammcj/tabload/api"
)

type TabLoad struct {
	window fyne.Window
	client *api.Client

	ready bool // Flag to indicate if the UI is fully Initialised

	// UI components
	adminKeyEntry           *widget.Entry
	apiURLEntry             *widget.Entry
	autosplitReserveCheck   *widget.Check
	autosplitReserveEntry   *widget.Entry
	cacheModeDropdown       *widget.Select
	cacheSizeCheck          *widget.Check
	cacheSizeEntry          *widget.Entry
	cancelDownloadButton    *widget.Button
	chunkSizeCheck          *widget.Check
	chunkSizeEntry          *widget.Entry
	connectButton           *widget.Button
	currentLorasLabel       *widget.Label
	currentModelInfo        *fyne.Container
	currentModelLabel       *widget.Label
	deletePresetButton      *widget.Button
	downloadButton          *widget.Button
	draftCacheModeDropdown  *widget.Select
	draftModelNameCheck     *widget.Check
	draftModelNameEntry     *widget.Entry
	draftRopeAlphaCheck     *widget.Check
	draftRopeAlphaEntry     *widget.Entry
	draftRopeScaleCheck     *widget.Check
	draftRopeScaleEntry     *widget.Entry
	excludeEntry            *widget.Entry
	fasttensorsCheck        *widget.Check
	folderNameEntry         *widget.Entry
	gpuSplitAutoCheck       *widget.Check
	gpuSplitCheck           *widget.Check
	gpuSplitEntry           *widget.Entry
	includeEntry            *widget.Entry
	loadLorasButton         *widget.Button
	loadModelButton         *widget.Button
	logsTextArea            *widget.Entry
	lorasDropdown           *widget.Select
	maxSeqLenCheck          *widget.Check
	maxSeqLenEntry          *widget.Entry
	modelsDropdown          *widget.Select
	numExpertsPerTokenCheck *widget.Check
	numExpertsPerTokenEntry *widget.Entry
	overrideBaseSeqLenCheck *widget.Check
	overrideBaseSeqLenEntry *widget.Entry
	presetDropdown          *widget.Select
	promptTemplateCheck     *widget.Check
	promptTemplateEntry     *widget.Entry
	repoIDEntry             *widget.Entry
	repoTypeDropdown        *widget.Select
	revisionEntry           *widget.Entry
	ropeAlphaCheck          *widget.Check
	ropeAlphaEntry          *widget.Entry
	ropeScaleCheck          *widget.Check
	ropeScaleEntry          *widget.Entry
	savePresetButton        *widget.Button
	tokenEntry              *widget.Entry
	unloadLorasButton       *widget.Button
	unloadModelButton       *widget.Button
	form                    *widget.Form
	connectionStatus        *widget.Label
	logsPane                *fyne.Container
	showingLogs             bool
	promptTemplateDropdown  *widget.Select // Templates
	temperatureSlider       *widget.Slider // Sampling parameters
	topKEntry               *widget.Entry  // Sampling parameters
	topPSlider              *widget.Slider // Sampling parameters
	minPSlider              *widget.Slider // Sampling parameters
	topASlider              *widget.Slider // Sampling parameters
	tfsSlider               *widget.Slider // Sampling parameters
	typicalPSlider          *widget.Slider // Sampling parameters
	repetitionPenaltySlider *widget.Slider // Sampling parameters
	presencePenaltySlider   *widget.Slider // Sampling parameters
	frequencyPenaltySlider  *widget.Slider // Sampling parameters
	mirostatModeSelect      *widget.Select // Sampling parameters
	mirostatTauSlider       *widget.Slider // Sampling parameters
	mirostatEtaSlider       *widget.Slider // Sampling parameters
	streamingCheck          *widget.Check  // Sampling parameters
	grammarEntry            *widget.Entry  // Sampling parameters
	logitBiasEntry          *widget.Entry  // Sampling parameters
	negativePromptEntry     *widget.Entry  // Sampling parameters
	jsonModeCheck           *widget.Check  // Sampling parameters
	speculativeNgramCheck   *widget.Check  // Sampling parameters
}

type Preset struct {
	Name               string   `json:"name"`
	MaxSeqLen          *int     `json:"max_seq_len,omitempty"`
	OverrideBaseSeqLen *int     `json:"override_base_seq_len,omitempty"`
	CacheSize          *int     `json:"cache_size,omitempty"`
	GPUSplitAuto       bool     `json:"gpu_split_auto,omitempty"`
	GPUSplit           string   `json:"gpu_split,omitempty"`
	RopeScale          *float64 `json:"rope_scale,omitempty"`
	RopeAlpha          *float64 `json:"rope_alpha,omitempty"`
	CacheMode          string   `json:"cache_mode,omitempty"`
	PromptTemplate     *string  `json:"prompt_template,omitempty"`
	NumExpertsPerToken *int     `json:"num_experts_per_token,omitempty"`
	DraftModelName     *string  `json:"draft_model_name,omitempty"`
	DraftRopeScale     *float64 `json:"draft_rope_scale,omitempty"`
	DraftRopeAlpha     *float64 `json:"draft_rope_alpha,omitempty"`
	DraftCacheMode     string   `json:"draft_cache_mode,omitempty"`
	Fasttensors        bool     `json:"fasttensors,omitempty"`
	AutosplitReserve   string   `json:"autosplit_reserve,omitempty"`
	ChunkSize          *int     `json:"chunk_size,omitempty"`
}
