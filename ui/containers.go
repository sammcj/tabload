package ui

import (
	"fmt"

	"fyne.io/fyne/v2/widget"
)

func (t *TabLoad) updateModelInfoContainer(data [][]string) {
	t.currentModelInfo.RemoveAll()
	for _, row := range data {
		if len(row) == 2 {
			label := widget.NewLabel(fmt.Sprintf("%s: %s", row[0], row[1]))
			t.currentModelInfo.Add(label)
		}
	}
	t.currentModelInfo.Refresh()
}

func (t *TabLoad) refreshCurrentLoras() {
	currentLoras, err := t.client.FetchCurrentLoras()
	if err != nil {
		fmt.Println("Error fetching current loras:", err)
		return
	}
	t.currentLorasLabel.SetText(currentLoras)
}
