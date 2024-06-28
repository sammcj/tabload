package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sammcj/tabload/logging"
)

func (t *TabLoad) getTemplatesDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		logging.Error("Failed to get user config directory", err)
		return "templates"
	}
	return filepath.Join(configDir, "tabload", "templates")
}

func (t *TabLoad) loadLocalTemplatesFromFiles() map[string]string {
	templatesDir := t.getTemplatesDir()
	templates := make(map[string]string)

	files, err := os.ReadDir(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			logging.Info("Templates directory does not exist, starting with empty templates")
			return templates
		}
		logging.Error("Failed to read templates directory", err)
		return templates
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}

		normalizedName := strings.TrimSuffix(file.Name(), ".txt")
		name := t.denormaliseTemplateName(normalizedName)

		content, err := os.ReadFile(filepath.Join(templatesDir, file.Name()))
		if err != nil {
			logging.Error(fmt.Sprintf("Failed to read template file %s", file.Name()), err)
			continue
		}

		templates[name] = string(content)
	}

	logging.Info(fmt.Sprintf("Loaded %d templates", len(templates)))
	return templates
}

func (t *TabLoad) showCreateTemplateDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Template Name")

	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetPlaceHolder("Enter template content")
	contentEntry.Resize(fyne.NewSize(500, 400))
	contentEntry.SetMinRowsVisible(5)

	content := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Name", nameEntry),
		),
		widget.NewLabel("Content:"),
		contentEntry,
	)

	dlg := dialog.NewCustom("Create New Template", "Save", content, t.window)
	dlg.SetOnClosed(func() {
		if nameEntry.Text != "" && contentEntry.Text != "" {
			t.saveNewTemplate(nameEntry.Text, contentEntry.Text)
		}
	})
	dlg.Resize(fyne.NewSize(500, 400))
	dlg.Show()
}

func (t *TabLoad) saveNewTemplate(name, content string) {
	templatesDir := t.getTemplatesDir()
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		logging.Error("Failed to create templates directory", err)
		dialog.ShowError(err, t.window)
		return
	}

	normalizedName := t.normaliseTemplateName(name)
	filePath := filepath.Join(templatesDir, normalizedName+".txt")

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		logging.Error("Failed to save new template", err)
		dialog.ShowError(err, t.window)
		return
	}

	t.refreshTemplateList()
	dialog.ShowInformation("Success", "Template saved successfully", t.window)
}

func (t *TabLoad) refreshTemplateList() {
	templates := t.loadTemplates()
	templateNames := make([]string, 0, len(templates)+1)
	for name := range templates {
		templateNames = append(templateNames, name)
	}
	templateNames = append(templateNames, "Create New...")

	if t.promptTemplateDropdown != nil {
		t.promptTemplateDropdown.Options = templateNames
		t.promptTemplateDropdown.Refresh()
	}
}

func (t *TabLoad) deleteTemplate(name string) {
	if strings.HasSuffix(name, " (server)") {
		dialog.ShowInformation("Cannot Delete", "Server-side templates cannot be deleted.", t.window)
		return
	}

	dialog.ShowConfirm("Delete Template", "Are you sure you want to delete this template?", func(confirm bool) {
		if confirm {
			normalizedName := t.normaliseTemplateName(name)
			filePath := filepath.Join(t.getTemplatesDir(), normalizedName+".txt")

			if err := os.Remove(filePath); err != nil {
				logging.Error("Failed to delete template file", err)
				dialog.ShowError(err, t.window)
				return
			}

			t.refreshTemplateList()
			dialog.ShowInformation("Success", "Template deleted successfully", t.window)
		}
	}, t.window)
}

func (t *TabLoad) normaliseTemplateName(name string) string {
	name = strings.ReplaceAll(name, " ", "__")
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' {
			return r
		}
		return '-'
	}, name)
}

func (t *TabLoad) denormaliseTemplateName(name string) string {
	name = strings.ReplaceAll(name, "__", " ")
	return strings.Map(func(r rune) rune {
		if r == '-' {
			return ' '
		}
		return r
	}, name)
}

func (t *TabLoad) loadTemplates() map[string]string {
	templates := make(map[string]string)

	// Load local templates
	localTemplates := t.loadLocalTemplatesFromFiles()
	for name, content := range localTemplates {
		templates[name] = content
	}

	// Fetch server-side templates
	if t.client != nil && t.client.BaseURL != "" {
		serverTemplates, err := t.client.FetchServerTemplates()
		if err != nil {
			logging.Error("Failed to fetch server templates", err)
		} else {
			for _, name := range serverTemplates {
				templates[name+" (server)"] = "" // We don't have the content, just the name
			}
		}
	} else {
		logging.Warn("Client not initialised or BaseURL not set, skipping server template fetch")
	}

	logging.Info(fmt.Sprintf("Loaded %d templates (%d local, %d server)", len(templates), len(localTemplates), len(templates)-len(localTemplates)))
	return templates
}
