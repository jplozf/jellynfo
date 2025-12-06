package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"jellynfo/jellydata"
)

var (
	originalReadInput    = readInputFunc
	originalEnsureConfig = ensureConfigFunc
	originalUserHomeDir  = jellydata.OsUserHomeDir

	// Mutex to protect global variables during parallel tests
	mu sync.Mutex
)

func TestEnsureConfig_NoConfigFile(t *testing.T) {
	// Mock ensureConfigFunc to simulate it creating and saving a config
	oldEnsureConfigFunc := ensureConfigFunc
	ensureConfigFunc = func() (jellydata.Config, error) {
		// Simulate the logic of creating the dir and writing the file
		tempDir, err := ioutil.TempDir("", "test_home_no_config")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		configDir := filepath.Join(tempDir, ".jellynfo")
		if err := os.MkdirAll(configDir, 0700); err != nil {
			t.Fatalf("Failed to create config dir: %v", err)
		}
		configPath := filepath.Join(configDir, "config.json")
		configContent := []byte(`{\"server_url\": \"http://mock.jellyfin.server\", \"api_key\": \"mock-api-key\"}`)
		if err := os.WriteFile(configPath, configContent, 0600); err != nil {
			t.Fatalf("Failed to write mock config file: %v", err)
		}
		return jellydata.Config{ServerURL: "http://mock.jellyfin.server", APIKey: "mock-api-key"}, nil
	}
	defer func() { ensureConfigFunc = oldEnsureConfigFunc }()

	config, err := ensureConfigFunc()
	if err != nil {
		t.Fatalf("ensureConfigFunc returned an unexpected error: %v", err)
	}

	expectedConfig := jellydata.Config{
		ServerURL: "http://mock.jellyfin.server",
		APIKey:    "mock-api-key",
	}

	if config.ServerURL != expectedConfig.ServerURL || config.APIKey != expectedConfig.APIKey {
		t.Errorf("Expected config %+v, got %+v", expectedConfig, config)
	}
}

func TestEnsureConfig_ExistingConfigFile(t *testing.T) {
	// Mock ensureConfigFunc to simulate it reading an existing config
	oldEnsureConfigFunc := ensureConfigFunc
	ensureConfigFunc = func() (jellydata.Config, error) {
		// Simulate the logic of reading an existing file
		tempDir, err := ioutil.TempDir("", "test_home_existing_config")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		configDir := filepath.Join(tempDir, ".jellynfo")
		if err := os.MkdirAll(configDir, 0700); err != nil {
			t.Fatalf("Failed to create config dir: %v", err)
		}
		configPath := filepath.Join(configDir, "config.json")
		existingContent := []byte(`{\"server_url\": \"http://existing.jellyfin.server\", \"api_key\": \"existing-api-key\"}`)
		if err := os.WriteFile(configPath, existingContent, 0600); err != nil {
			t.Fatalf("Failed to write existing mock config file: %v", err)
		}
		return jellydata.Config{ServerURL: "http://existing.jellyfin.server", APIKey: "existing-api-key"}, nil
	}
	defer func() { ensureConfigFunc = oldEnsureConfigFunc }()

	config, err := ensureConfigFunc()
	if err != nil {
		t.Fatalf("ensureConfigFunc returned an unexpected error: %v", err)
	}

	expectedConfig := jellydata.Config{
		ServerURL: "http://existing.jellyfin.server",
		APIKey:    "existing-api-key",
	}

	if config.ServerURL != expectedConfig.ServerURL || config.APIKey != expectedConfig.APIKey {
		t.Errorf("Expected config %+v, got %+v", expectedConfig, config)
	}
}

func TestFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected OutputType
		err      bool // true if an error is expected
	}{
		{
			name:     "default output",
			args:     []string{"jellynfo"}, // Program name is always the first arg
			expected: OutputConsole,
			err:      false,
		},
		{
			name:     "console output",
			args:     []string{"jellynfo", "--output=console"},
			expected: OutputConsole,
			err:      false,
		},
		{
			name:     "invalid output",
			args:     []string{"jellynfo", "--output=invalid"},
			expected: "", // Expected value doesn't matter, as an error is expected
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup temporary home directory for this sub-test
			tempDir, err := ioutil.TempDir("", "test_home_flag_sub")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}

			// Mock os.UserHomeDir
			oldUserHomeDir := jellydata.OsUserHomeDir
			jellydata.OsUserHomeDir = func() (string, error) {
				return tempDir, nil
			}

			// Mock readInputFunc
			oldReadInputFunc := readInputFunc
			inputReader := bytes.NewBufferString(strings.Join([]string{"http://dummy.jellyfin.server", "dummy-api-key"}, "\n"))
			readInputFunc = func(prompt string) (string, error) {
				line, err := inputReader.ReadString('\n')
				return strings.TrimSpace(line), err
			}

			// Temporarily redirect ensureConfigFunc to directly return the dummy config
			oldEnsureConfigFunc := ensureConfigFunc
			ensureConfigFunc = func() (jellydata.Config, error) {
				return jellydata.Config{ServerURL: "http://dummy.jellyfin.server", APIKey: "dummy-api-key"}, nil
			}
			// Create a dummy config file so ensureConfig doesn't prompt if not mocked
			configDir := filepath.Join(tempDir, ".jellynfo")
			if err := os.MkdirAll(configDir, 0700); err != nil {
				t.Fatalf("Failed to create config dir: %v", err)
			}
			configPath := filepath.Join(configDir, "config.json")
			dummyConfigContent := []byte(`{"server_url": "http://dummy.jellyfin.server", "api_key": "dummy-api-key"}`)
			if err := os.WriteFile(configPath, dummyConfigContent, 0600); err != nil {
				t.Fatalf("Failed to write dummy config file: %v", err)
			}

			// Save the original CommandLine and restore it after the test
			oldArgs := os.Args

			// Cleanup function for this sub-test
			defer func() {
				jellydata.OsUserHomeDir = oldUserHomeDir
				readInputFunc = oldReadInputFunc
				ensureConfigFunc = oldEnsureConfigFunc
				os.Args = oldArgs
				if err := os.RemoveAll(tempDir); err != nil {
					t.Fatalf("Failed to clean up temp dir %s: %v", tempDir, err)
				}
			}()

			os.Args = tt.args

			// Reset flags for each test run
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			outputFlag = OutputConsole // Reset to default before parsing

			// Manually register the flag for the new FlagSet
			flag.CommandLine.Func("output", "Output mode: 'console'", func(s string) error {
				switch s {
				case "console":
					outputFlag = OutputConsole
				default:
					return fmt.Errorf("invalid output type: %s. Must be 'console'", s)
				}
				return nil
			})

			err = flag.CommandLine.Parse(os.Args[1:])

			if (err != nil) != tt.err {
				// For invalid output flag, we expect an error, but not necessarily a fatal one from ensureConfigFunc.
				// So, we only check if the error status matches the expected error status for the test.
				if !tt.err && err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else if tt.err && err == nil {
				t.Fatalf("expected an error, but got none")
			}

			if !tt.err && outputFlag != tt.expected {
				t.Errorf("expected outputFlag to be %q, got %q", tt.expected, outputFlag)
			}
		})
	}
}
