package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"jellynfo/jellydata"
	"jellynfo/output"
)

type OutputType string

const (
	OutputConsole OutputType = "console"
	OutputUI      OutputType = "ui"
)

var outputFlag OutputType

func init() {
	flag.Func("output", "Output mode: 'console' (default) or 'ui'", func(s string) error {
		switch s {
		case "console":
			outputFlag = OutputConsole
		case "ui":
			outputFlag = OutputUI
		default:
			return fmt.Errorf("invalid output type: %s. Must be 'console' or 'ui'", s)
		}
		return nil
	})
	outputFlag = OutputConsole
}

type ReadInputFunc func(prompt string) (string, error)

var readInputFunc ReadInputFunc = defaultReadInput

func defaultReadInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(input), nil
}

type EnsureConfigFunc func() (jellydata.Config, error)

var ensureConfigFunc EnsureConfigFunc = defaultEnsureConfig

func defaultEnsureConfig() (jellydata.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return jellydata.Config{}, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".jellynfo")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		log.Printf("Creating config directory: %s", configDir)
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return jellydata.Config{}, fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	configPath := filepath.Join(configDir, "config.json")

	var config jellydata.Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found: %s. Prompting for credentials...", configPath)

		serverURL, err := readInputFunc("Enter Jellyfin server address (e.g., http://localhost:8096): ")
		if err != nil {
			return jellydata.Config{}, err
		}

		apiKey, err := readInputFunc("Enter Jellyfin API token: ")
		if err != nil {
			return jellydata.Config{}, err
		}

		config = jellydata.Config{
			ServerURL: serverURL,
			APIKey:    apiKey,
		}

		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return jellydata.Config{}, fmt.Errorf("failed to marshal config: %w", err)
		}

		if err := os.WriteFile(configPath, configBytes, 0600); err != nil {
			return jellydata.Config{}, fmt.Errorf("failed to write config file: %w", err)
		}
		log.Printf("Config saved to: %s", configPath)
	} else if err == nil {
		// Config file exists, read it
		configFile, err := os.Open(configPath)
		if err != nil {
			return jellydata.Config{}, fmt.Errorf("failed to open config file: %w", err)
		}
		defer configFile.Close()

		jsonParser := json.NewDecoder(configFile)
		if err = jsonParser.Decode(&config); err != nil {
			return jellydata.Config{}, fmt.Errorf("failed to parse config file: %w", err)
		}
	} else {
		return jellydata.Config{}, fmt.Errorf("failed to stat config file: %w", err)
	}

	return config, nil
}

func main() {
	flag.Parse()

	config, err := ensureConfigFunc()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	ctx := context.Background()
	data, err := jellydata.GetJellyfinData(ctx, config)
	if err != nil {
		log.Fatalf("Error fetching Jellyfin data: %v", err)
	}

	var outputter output.Outputter

	switch outputFlag {
	case OutputConsole:
		outputter = output.NewConsoleOutputter()
	case OutputUI:
		outputter = output.NewFyneGUIOutputter()
	}

	if err := outputter.Display(data); err != nil {
		log.Fatalf("Error displaying Jellyfin data: %v", err)
	}
}
