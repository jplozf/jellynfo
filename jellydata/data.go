package jellydata

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sj14/jellyfin-go/api"
)

// Config represents the structure of the config.json file
type Config struct {
	ServerURL string `json:"server_url"`
	APIKey    string `json:"api_key"`
}

// SystemInfo holds the relevant system information
type SystemInfo struct {
	ServerName     string
	Version        string
	OperatingSystem string
}

// SessionInfo holds the relevant session information
type SessionInfo struct {
	User           string
	Device         string
	Client         string
	IPAddress      string
	LastActivity   time.Time
	NowPlayingItem string
}

// JellyfinData aggregates all fetched information
type JellyfinData struct {
	System          SystemInfo
	Sessions        []SessionInfo
	PlayingSessions int32 // New field for count of playing sessions
}

// GetJellyfinDataFunc defines the function signature for fetching Jellyfin data.
type GetJellyfinDataFunc func(ctx context.Context, config Config) (*JellyfinData, error)

// DefaultGetJellyfinData is the actual implementation that fetches data from the real API.
var GetJellyfinData = defaultGetJellyfinData

var OsUserHomeDir = os.UserHomeDir

func defaultGetJellyfinData(ctx context.Context, config Config) (*JellyfinData, error) {
	// Add default protocol scheme if missing
	if !strings.HasPrefix(config.ServerURL, "http://") && !strings.HasPrefix(config.ServerURL, "https://") {
		config.ServerURL = "http://" + config.ServerURL
	}

	// Configure API client
	configuration := api.NewConfiguration()
	configuration.Servers = api.ServerConfigurations{
		{
			URL: config.ServerURL,
		},
	}
	configuration.AddDefaultHeader("X-Emby-Token", config.APIKey) // Authenticate with API key

	// Optional: Set a custom HTTP client with a timeout
	configuration.HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	apiClient := api.NewAPIClient(configuration)

	data := &JellyfinData{}

	// Get system information
	systemInfo, resp, err := apiClient.SystemAPI.GetSystemInfo(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting system info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get system info. Status: %s", resp.Status)
	}

	if serverName := systemInfo.ServerName.Get(); serverName != nil {
		data.System.ServerName = *serverName
	}
	if version := systemInfo.Version.Get(); version != nil {
		data.System.Version = *version
	}
	if osName := systemInfo.OperatingSystem.Get(); osName != nil && *osName != "" {
		data.System.OperatingSystem = *osName
	} else {
		data.System.OperatingSystem = "(not available)"
	}

	// Get sessions information
	sessions, _, err := apiClient.SessionAPI.GetSessions(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting sessions: %v", err)
	}

	playingSessionsCount := int32(0)
	if len(sessions) == 0 {
		data.Sessions = []SessionInfo{} // Ensure it's an empty slice, not nil
	} else {
		for _, session := range sessions {
			sInfo := SessionInfo{}
			if userName := session.UserName.Get(); userName != nil {
				sInfo.User = *userName
			}
			if deviceName := session.DeviceName.Get(); deviceName != nil {
				sInfo.Device = *deviceName
			}
			if client := session.Client.Get(); client != nil {
				sInfo.Client = *client
			}
			if ipAddress := session.RemoteEndPoint.Get(); ipAddress != nil {
				sInfo.IPAddress = *ipAddress
			}
			if lastActivityDate := session.LastActivityDate; lastActivityDate != nil {
				sInfo.LastActivity = *lastActivityDate
			}

			if nowPlayingItem := session.NowPlayingItem.Get(); nowPlayingItem != nil {
				mediaType := "Unknown"
				if nowPlayingItem.Type != nil {
					mediaType = string(*nowPlayingItem.Type)
				}
				if itemName := nowPlayingItem.Name.Get(); itemName != nil {
					if mediaType == "Episode" {
						if seriesName := nowPlayingItem.SeriesName.Get(); seriesName != nil {
							sInfo.NowPlayingItem = fmt.Sprintf("%s - %s (%s)", *seriesName, *itemName, mediaType)
						} else {
							sInfo.NowPlayingItem = fmt.Sprintf("%s (%s)", *itemName, mediaType)
						}
					} else {
						sInfo.NowPlayingItem = fmt.Sprintf("%s (%s)", *itemName, mediaType)
					}
					playingSessionsCount++
				} else {
					sInfo.NowPlayingItem = "Information unavailable"
				}
			} else {
				sInfo.NowPlayingItem = "Nothing"
			}
			data.Sessions = append(data.Sessions, sInfo)
		}
	}
	data.PlayingSessions = playingSessionsCount

	return data, nil
}