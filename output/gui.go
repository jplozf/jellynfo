package output

import (
	"fmt"
	"jellynfo/jellydata"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// FyneGUIOutputter implements the Outputter interface for Fyne GUI.
type FyneGUIOutputter struct{}

// NewFyneGUIOutputter creates a new FyneGUIOutputter.
func NewFyneGUIOutputter() *FyneGUIOutputter {
	return &FyneGUIOutputter{}
}

// Display creates and shows the Fyne GUI with Jellyfin data.
func (f *FyneGUIOutputter) Display(data *jellydata.JellyfinData) error {
	a := app.New()
	w := a.NewWindow("Jellyfin Information")

	// System Info
	var systemInfoObjects []fyne.CanvasObject
	systemInfoObjects = append(systemInfoObjects, widget.NewLabelWithStyle("System Information:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	systemInfoObjects = append(systemInfoObjects, widget.NewLabel(fmt.Sprintf("Server Name: %s", data.System.ServerName)))
	systemInfoObjects = append(systemInfoObjects, widget.NewLabel(fmt.Sprintf("Version: %s", data.System.Version)))
	systemInfoObjects = append(systemInfoObjects, widget.NewLabel(fmt.Sprintf("Operating System: %s", data.System.OperatingSystem)))
	systemInfoCard := widget.NewCard("", "", container.NewVBox(systemInfoObjects...))

	// Sessions Info
	var sessionLabels []fyne.CanvasObject

	sessionLabels = append(sessionLabels, widget.NewLabelWithStyle("Connected Users:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

	if len(data.Sessions) == 0 {
		sessionLabels = append(sessionLabels, widget.NewLabel("No active sessions found."))
	} else {
		sessionLabels = append(sessionLabels, widget.NewLabel(fmt.Sprintf("Number of active sessions: %d", len(data.Sessions))))
		sessionLabels = append(sessionLabels, widget.NewLabel(fmt.Sprintf("Number of playing sessions: %d", data.PlayingSessions)))

		headers := []string{"User", "Device", "Client", "IP Address", "Last Activity", "Now Playing"}
		tableData := make([][]string, len(data.Sessions))

		for i, session := range data.Sessions {
			tableData[i] = []string{
				session.User,
				session.Device,
				session.Client,
				session.IPAddress,
				session.LastActivity.Format("2006-01-02 15:04:05"),
				session.NowPlayingItem,
			}
		}

		sessionsTable := widget.NewTable(
			func() (int, int) { return len(tableData), len(headers) },
			func() fyne.CanvasObject {
				return widget.NewLabel("Wide content")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				label := o.(*widget.Label)
				if i.Row == 0 {
					label.SetText(headers[i.Col])
					label.TextStyle = fyne.TextStyle{Bold: true}
				} else {
					label.SetText(tableData[i.Row-1][i.Col])
					label.TextStyle = fyne.TextStyle{Bold: false}
				}
			},
		)

		// Set headers manually (first row)
		sessionsTable.Length = func() (int, int) { return len(tableData) + 1, len(headers) }

		sessionsTable.SetColumnWidth(0, 100)
		sessionsTable.SetColumnWidth(1, 150)
		sessionsTable.SetColumnWidth(2, 150)
		sessionsTable.SetColumnWidth(3, 120)
		sessionsTable.SetColumnWidth(4, 180)
		sessionsTable.SetColumnWidth(5, 250)

		sessionLabels = append(sessionLabels, sessionsTable)
	}

	sessionsCard := widget.NewCard("", "", container.NewVBox(sessionLabels...))

	content := container.NewVBox(
		systemInfoCard,
		sessionsCard,
	)

	w.SetContent(container.NewScroll(content))
	w.Resize(fyne.NewSize(600, 400)) // Set a default window size
	w.ShowAndRun()
	return nil
}