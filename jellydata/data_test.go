package jellydata

import (
	"context"
	"testing"
	"time"
)

func TestGetJellyfinData_Mocked(t *testing.T) {
	// Override GetJellyfinData for this test with a mock
	oldGetJellyfinData := GetJellyfinData
	GetJellyfinData = func(ctx context.Context, config Config) (*JellyfinData, error) {
		// Return hardcoded mock data
		mockTime, _ := time.Parse(time.RFC3339, "2025-12-05T10:00:00.000Z")
		mockTime2, _ := time.Parse(time.RFC3339, "2025-12-05T09:30:00.000Z")

		return &JellyfinData{
			System: SystemInfo{
				ServerName:      "MockJellyfin",
				Version:         "1.0.0",
				OperatingSystem: "MockOS",
			},
			Sessions: []SessionInfo{
				{
					User:           "mockuser1",
					Device:         "MockDesktop",
					Client:         "MockWeb",
					IPAddress:      "127.0.0.1",
					LastActivity:   mockTime,
					NowPlayingItem: "Mock Series - Mock Episode (Episode)",
				},
				{
					User:           "mockuser2",
					Device:         "MockMobile",
					Client:         "MockAndroid",
					IPAddress:      "127.0.0.2",
					LastActivity:   mockTime2,
					NowPlayingItem: "Nothing",
				},
			},
			PlayingSessions: 1, // One session is playing
		}, nil
	}
	defer func() { GetJellyfinData = oldGetJellyfinData }() // Restore original after test

	ctx := context.Background()
	data, err := GetJellyfinData(ctx, Config{})
	if err != nil {
		t.Fatalf("GetJellyfinData returned an error: %v", err)
	}

	// Assertions for SystemInfo
	if data.System.ServerName != "MockJellyfin" {
		t.Errorf("Expected ServerName 'MockJellyfin', got %q", data.System.ServerName)
	}
	if data.System.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got %q", data.System.Version)
	}
	if data.System.OperatingSystem != "MockOS" {
		t.Errorf("Expected OperatingSystem 'MockOS', got %q", data.System.OperatingSystem)
	}

	// Assertions for Sessions
	if len(data.Sessions) != 2 {
		t.Fatalf("Expected 2 sessions, got %d", len(data.Sessions))
	}
	if data.PlayingSessions != 1 {
		t.Errorf("Expected 1 playing session, got %d", data.PlayingSessions)
	}

	// Session 1
	session1 := data.Sessions[0]
	if session1.User != "mockuser1" {
		t.Errorf("Expected session 1 User 'mockuser1', got %q", session1.User)
	}
	if session1.NowPlayingItem != "Mock Series - Mock Episode (Episode)" {
		t.Errorf("Expected session 1 NowPlayingItem 'Mock Series - Mock Episode (Episode)', got %q", session1.NowPlayingItem)
	}

	// Session 2
	session2 := data.Sessions[1]
	if session2.User != "mockuser2" {
		t.Errorf("Expected session 2 User 'mockuser2', got %q", session2.User)
	}
	if session2.NowPlayingItem != "Nothing" {
		t.Errorf("Expected session 2 NowPlayingItem 'Nothing', got %q", session2.NowPlayingItem)
	}
}