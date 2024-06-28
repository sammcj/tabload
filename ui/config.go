package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sammcj/tabload/logging"
	"github.com/spf13/viper"
)

type Config struct {
	AutoConnect         bool        `json:"auto_connect"`
	LastConnectedServer string      `json:"last_connected_server"`
	APIURL              string      `json:"api_url"`
	DefaultParams       ModelParams `json:"default_params"`
}

type ModelParams struct {
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

var (
	homeDir, _      = os.UserHomeDir()
	configPath      = filepath.Join(homeDir, ".config", "tabload") // $HOME/.config/tabload
	configFile      = filepath.Join(configPath, "config.json")
	presetsFilePath = filepath.Join(configPath, "presets.json") // $HOME/.config/tabload/presets.json
	config          Config
)

func initConfig() {
	logging.Info(fmt.Sprintf("Attempting to read config from: %s", configFile))
	viper.SetConfigFile(configFile)

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			logging.Warn(fmt.Sprintf("Config file not found at %s. Using defaults.", configFile))
			config = Config{
				AutoConnect:         false,
				LastConnectedServer: "",
				APIURL:              "http://localhost:5000",
				DefaultParams:       ModelParams{},
			}
			// Write default config
			if err := saveConfig(); err != nil {
				logging.Error("Failed to write default config", err)
			}
		} else {
			logging.Error("Failed to read config file: %v", err)
		}
	} else {
		if err := json.Unmarshal(data, &config); err != nil {
			logging.Error("Failed to unmarshal config: %v", err)
			// Use default config if unmarshal fails
			config = Config{
				AutoConnect:         false,
				LastConnectedServer: "",
				APIURL:              "http://localhost:5000",
				DefaultParams:       ModelParams{},
			}
		} else {
			logging.Info("Successfully read config file")
		}
	}

	logging.Info(fmt.Sprintf("Loaded config: %+v", config))
}

func saveConfig() error {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "tabload")
	configFile := filepath.Join(configPath, "config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (t *TabLoad) SaveDefaultParams() error {
	params := t.createPresetFromFields()
	config.DefaultParams = ModelParams{
		MaxSeqLen:          params.MaxSeqLen,
		OverrideBaseSeqLen: params.OverrideBaseSeqLen,
		CacheSize:          params.CacheSize,
		GPUSplitAuto:       params.GPUSplitAuto,
		GPUSplit:           params.GPUSplit,
		RopeScale:          params.RopeScale,
		RopeAlpha:          params.RopeAlpha,
		CacheMode:          params.CacheMode,
		PromptTemplate:     params.PromptTemplate,
		NumExpertsPerToken: params.NumExpertsPerToken,
		DraftModelName:     params.DraftModelName,
		DraftRopeScale:     params.DraftRopeScale,
		DraftRopeAlpha:     params.DraftRopeAlpha,
		DraftCacheMode:     params.DraftCacheMode,
		Fasttensors:        params.Fasttensors,
		AutosplitReserve:   params.AutosplitReserve,
		ChunkSize:          params.ChunkSize,
	}
	viper.Set("default_params", config.DefaultParams)
	return viper.WriteConfig()
}

func (t *TabLoad) LoadDefaultParams() {
	preset := Preset{
		MaxSeqLen:          config.DefaultParams.MaxSeqLen,
		OverrideBaseSeqLen: config.DefaultParams.OverrideBaseSeqLen,
		CacheSize:          config.DefaultParams.CacheSize,
		GPUSplitAuto:       config.DefaultParams.GPUSplitAuto,
		GPUSplit:           config.DefaultParams.GPUSplit,
		RopeScale:          config.DefaultParams.RopeScale,
		RopeAlpha:          config.DefaultParams.RopeAlpha,
		CacheMode:          config.DefaultParams.CacheMode,
		PromptTemplate:     config.DefaultParams.PromptTemplate,
		NumExpertsPerToken: config.DefaultParams.NumExpertsPerToken,
		DraftModelName:     config.DefaultParams.DraftModelName,
		DraftRopeScale:     config.DefaultParams.DraftRopeScale,
		DraftRopeAlpha:     config.DefaultParams.DraftRopeAlpha,
		DraftCacheMode:     config.DefaultParams.DraftCacheMode,
		Fasttensors:        config.DefaultParams.Fasttensors,
		AutosplitReserve:   config.DefaultParams.AutosplitReserve,
		ChunkSize:          config.DefaultParams.ChunkSize,
	}
	t.applyPresetToFields(&preset)
}
func (t *TabLoad) SaveConfig() error {
	return viper.WriteConfig()
}

func (t *TabLoad) GetConfig() Config {
	return config
}

func (t *TabLoad) SetConfig(newConfig Config) {
	config = newConfig
	viper.Set("auto_connect", config.AutoConnect)
	viper.Set("last_connected_server", config.LastConnectedServer)
	viper.Set("api_url", config.APIURL)
	viper.Set("default_params", config.DefaultParams)
}
func (t *TabLoad) saveLastConnectedServer(url string) error {
	config.LastConnectedServer = url
	viper.Set("last_connected_server", url)
	return viper.WriteConfig()
}

func (t *TabLoad) loadLastConnectedServer() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(configDir, "tabload", "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return "", err
	}

	return config.LastConnectedServer, nil
}

func (t *TabLoad) saveAutoConnectSetting(autoConnect bool) {
	viper.Set("auto_connect", autoConnect)
	if err := viper.WriteConfig(); err != nil {
		logging.Error("Error saving auto-connect setting", err)
		return
	}

	if autoConnect {
		viper.Set("last_connected_server", t.apiURLEntry.Text)
		if err := viper.WriteConfig(); err != nil {
			logging.Error("Error saving last connected server", err)
			return
		}
	}

	logging.Info(fmt.Sprintf("Auto-connect setting saved: %v", autoConnect))
}
