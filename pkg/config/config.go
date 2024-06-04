package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const CHANGESET_DIRECTORY string = ".changeset"
const CONFIG_FILENAME string = "config.json"

type Plugin struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	SHA256        string `json:"sha256"`
	VersionedFile string `json:"versionedFile"`
}

type Config struct {
	Plugin Plugin `json:"plugin"`
}

func GetConfig() (Config, error) {
	filepath := filepath.Join(CHANGESET_DIRECTORY, CONFIG_FILENAME)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return Config{}, errors.New("config file not found")
	}
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
