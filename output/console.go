package output

import (
	"fmt"
	"jellynfo/jellydata"
	"strings"
)

// ConsoleOutputter implements the Outputter interface for console output.
type ConsoleOutputter struct{}

// NewConsoleOutputter creates a new ConsoleOutputter.
func NewConsoleOutputter() *ConsoleOutputter {
	return &ConsoleOutputter{}
}

// Display prints the JellyfinData to the console.
func (c *ConsoleOutputter) Display(data *jellydata.JellyfinData) error {
	fmt.Println("Successfully connected to Jellyfin server.")
	fmt.Printf("Server Name: %s\n", data.System.ServerName)
	fmt.Printf("Version: %s\n", data.System.Version)
	fmt.Printf("Operating System: %s\n", data.System.OperatingSystem)

	fmt.Println("\n--- Connected Users ---")

	if len(data.Sessions) == 0 {
		fmt.Println("No active sessions found.")
		return nil
	}

	fmt.Printf("Number of active  sessions: %d\n", len(data.Sessions))
	fmt.Printf("Number of playing sessions: %d\n\n", data.PlayingSessions)

	// Print table header
	fmt.Printf("%-15s %-25s %-25s %-15s %-25s %s\n", "User", "Device", "Client", "IP Address", "Last Activity", "Now Playing")
	fmt.Println(strings.Repeat("-", 130)) // Adjust width as needed

	// Print session data
	for _, session := range data.Sessions {
		fmt.Printf("%-15s %-25s %-25s %-15s %-25s %s\n",
			session.User,
			session.Device,
			session.Client,
			session.IPAddress,
			session.LastActivity.Format("2006-01-02 15:04:05"),
			session.NowPlayingItem,
		)
	}
	fmt.Println("")
	return nil
}
