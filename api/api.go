package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sammcj/tabload/logging"
)

func NewClient(baseURL, adminKey string) *Client {
	return &Client{
		BaseURL:  baseURL,
		AdminKey: adminKey,
	}
}

func (c *Client) makeHTTPRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	url := c.BaseURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("X-api-key", c.AdminKey)
	if method == http.MethodPost {
		req.Header.Add("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) FetchModels() ([]string, error) {
	body, err := c.makeHTTPRequest(http.MethodGet, "/v1/model/list", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching models: %w", err)
	}

	var response struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	models := make([]string, len(response.Data))
	for i, model := range response.Data {
		models[i] = model.ID
	}

	return models, nil
}

func (c *Client) FetchDraftModels() ([]string, error) {
	body, err := c.makeHTTPRequest(http.MethodGet, "/v1/model/draft/list", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching draft models: %w", err)
	}

	var response struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	draftModels := make([]string, len(response.Data))
	for i, model := range response.Data {
		draftModels[i] = model.ID
	}

	return draftModels, nil
}

func (c *Client) FetchLoras() ([]string, error) {
	body, err := c.makeHTTPRequest(http.MethodGet, "/v1/lora/list", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching loras: %w", err)
	}

	var response struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	loras := make([]string, len(response.Data))
	for i, lora := range response.Data {
		loras[i] = lora.ID
	}

	return loras, nil
}

func (c *Client) FetchTemplates() ([]string, error) {
	url := c.BaseURL + "/v1/template/list"
	logging.Info(fmt.Sprintf("Fetching templates from: %s", url))

	body, err := c.makeHTTPRequest(http.MethodGet, "/v1/template/list", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching templates: %w", err)
	}

	var response struct {
		Data []string `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return response.Data, nil
}

func (c *Client) FetchOverrides() ([]string, error) {
	body, err := c.makeHTTPRequest(http.MethodGet, "/v1/sampling/override/list", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching overrides: %w", err)
	}

	var response struct {
		Presets []string `json:"presets"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return response.Presets, nil
}

type Model struct {
	ID         string
	Parameters struct {
		MaxSeqLen int
		CacheSize int
		RopeScale float64
		RopeAlpha float64
		Draft     *struct {
			ID         string
			Parameters struct {
				RopeScale float64
				RopeAlpha float64
			}
		}
	}
}

func (c *Client) FetchCurrentModel() (*Model, error) {
	body, err := c.makeHTTPRequest(http.MethodGet, "/v1/model", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching current model: %w", err)
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
	if err := json.Unmarshal(body, &response); err != nil {
		return &Model{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	model := fmt.Sprintf("%s (context: %d, cache size: %d, rope scale: %.2f, rope alpha: %.2f)",
		response.ID, response.Parameters.MaxSeqLen, response.Parameters.CacheSize,
		response.Parameters.RopeScale, response.Parameters.RopeAlpha)

	if response.Parameters.Draft.ID != "" {
		model += fmt.Sprintf(" | %s (rope scale: %.2f, rope alpha: %.2f)",
			response.Parameters.Draft.ID, response.Parameters.Draft.Parameters.RopeScale,
			response.Parameters.Draft.Parameters.RopeAlpha)
	}

	return &Model{
		ID: response.ID,
		Parameters: struct {
			MaxSeqLen int
			CacheSize int
			RopeScale float64
			RopeAlpha float64
			Draft     *struct {
				ID         string
				Parameters struct {
					RopeScale float64
					RopeAlpha float64
				}
			}
		}{
			MaxSeqLen: response.Parameters.MaxSeqLen,
			CacheSize: response.Parameters.CacheSize,
		},
	}, nil
}

func (c *Client) FetchCurrentLoras() (string, error) {
	body, err := c.makeHTTPRequest(http.MethodGet, "/v1/lora", nil)
	if err != nil {
		return "", fmt.Errorf("fetching current loras: %w", err)
	}

	var response struct {
		Data []struct {
			ID      string  `json:"id"`
			Scaling float64 `json:"scaling"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("unmarshalling response: %w", err)
	}

	loras := make([]string, len(response.Data))
	for i, lora := range response.Data {
		loras[i] = fmt.Sprintf("%s (scaling: %.2f)", lora.ID, lora.Scaling)
	}

	return strings.Join(loras, ", "), nil
}

func (c *Client) LoadModel(modelName string, params map[string]interface{}) error {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshalling params: %w", err)
	}

	req, err := c.makeHTTPRequest(http.MethodPost, "/v1/model/load", strings.NewReader(string(jsonData)))
	logging.Debug("Request: " + string(jsonData))
	if err != nil {
		logging.Debug("Error loading model: " + err.Error())
		logging.Debug("Response: " + string(req))
		return fmt.Errorf("loading model: %w", err)
	}

	return nil
}

func (c *Client) LoadLoras(loras []string, scalings []float64) error {
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
		return fmt.Errorf("marshalling request: %w", err)
	}

	_, err = c.makeHTTPRequest(http.MethodPost, "/v1/lora/load", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("loading loras: %w", err)
	}

	return nil
}

func (c *Client) UnloadModel() error {
	_, err := c.makeHTTPRequest(http.MethodPost, "/v1/model/unload", nil)
	if err != nil {
		return fmt.Errorf("unloading model: %w", err)
	}
	return nil
}

func (c *Client) UnloadLoras() error {
	_, err := c.makeHTTPRequest(http.MethodPost, "/v1/lora/unload", nil)
	if err != nil {
		return fmt.Errorf("unloading loras: %w", err)
	}
	return nil
}

func (c *Client) LoadTemplate(promptTemplate string) error {
	request := map[string]string{
		"name": promptTemplate,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshalling request: %w", err)
	}

	_, err = c.makeHTTPRequest(http.MethodPost, "/v1/template/switch", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("loading template: %w", err)
	}

	return nil
}

func (c *Client) UnloadTemplate() error {
	_, err := c.makeHTTPRequest(http.MethodPost, "/v1/template/unload", nil)
	if err != nil {
		return fmt.Errorf("unloading template: %w", err)
	}
	return nil
}

func (c *Client) LoadOverride(samplerOverride string) error {
	request := map[string]string{
		"preset": samplerOverride,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshalling request: %w", err)
	}

	_, err = c.makeHTTPRequest(http.MethodPost, "/v1/sampling/override/switch", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("loading override: %w", err)
	}

	return nil
}

func (c *Client) UnloadOverride() error {
	_, err := c.makeHTTPRequest(http.MethodPost, "/v1/sampling/override/unload", nil)
	if err != nil {
		return fmt.Errorf("unloading override: %w", err)
	}
	return nil
}

func (c *Client) Download(params map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("marshalling params: %w", err)
	}

	body, err := c.makeHTTPRequest(http.MethodPost, "/v1/download", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("downloading: %w", err)
	}

	var response struct {
		DownloadPath string `json:"download_path"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("unmarshalling response: %w", err)
	}

	return response.DownloadPath, nil
}

func (c *Client) CancelDownload() error {
	_, err := c.makeHTTPRequest(http.MethodPost, "/v1/download/cancel", nil)
	if err != nil {
		return fmt.Errorf("cancelling download: %w", err)
	}
	return nil
}

func (c *Client) SaveTemplate(name, content string) error {
	params := map[string]string{
		"name":    name,
		"content": content,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshalling params: %w", err)
	}

	_, err = c.makeHTTPRequest(http.MethodPost, "/v1/template/save", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("saving template: %w", err)
	}

	return nil
}

func (c *Client) FetchServerTemplates() ([]string, error) {
	body, err := c.makeHTTPRequest(http.MethodGet, "/v1/template/list", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching templates: %w", err)
	}

	var response struct {
		Data []string `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return response.Data, nil
}
