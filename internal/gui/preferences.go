package gui

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UserPreferences struct {
	TutorialCompleted bool `json:"tutorial_completed"`
	FirstLaunch       bool `json:"first_launch"`
}

func GetPreferencesPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	configDir := filepath.Join(home, ".config", "systemdesignsim")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	
	return filepath.Join(configDir, "preferences.json"), nil
}

func LoadPreferences() (*UserPreferences, error) {
	path, err := GetPreferencesPath()
	if err != nil {
		return &UserPreferences{FirstLaunch: true}, nil
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &UserPreferences{FirstLaunch: true}, nil
		}
		return nil, err
	}
	
	var prefs UserPreferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return &UserPreferences{FirstLaunch: true}, nil
	}
	
	return &prefs, nil
}

func SavePreferences(prefs *UserPreferences) error {
	path, err := GetPreferencesPath()
	if err != nil {
		return err
	}
	
	data, err := json.Marshal(prefs)
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}
