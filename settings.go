// This package was made, to easily get needed settings from a file.
// Supported file types are: json and yaml.
//
// This package uses https://github.com/spf13/viper: Copyright © 2014 Steve Francia <spf@spf13.com>.
package settings

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type fsNotify struct {
	watcher *fsnotify.Watcher
	error   error
}

type Settings struct {
	Data      *viper.Viper
	Error     error
	content   string
	fileNames []string
	mux       sync.Mutex
}

// New initializes settings from a file or from multiple files under given directory.
func New(settingsFile string) *Settings {
	s := &Settings{}
	s.Data = viper.New()
	return s.load(settingsFile)
}

// NewFromContent initializes settings from a given content.
func NewFromContent(content string) *Settings {
	s := &Settings{content: content}
	s.Data = viper.New()
	ext := getExtensionByContent(content)
	if ext == "unsupported" {
		return &Settings{Error: fmt.Errorf("settings.NewFromContent :: unsupported content type")}
	}
	s.Data.SetConfigType(ext)
	_ = s.Data.ReadConfig(bytes.NewBuffer([]byte(content)))
	return s
}

// Merge merges initialized settings with a given file or directory.
func (s *Settings) Merge(settingsFile string) *Settings {
	return s.load(settingsFile)
}

// GetAllKeys returns all keys holding a value, regardless of where they are set.
// Nested keys are returned with a v.key delimiter separator
func (s *Settings) GetAllKeys() ([]string, error) {
	if s.Error != nil {
		return nil, fmt.Errorf("settings.GetAllKeys :: %s", s.Error)
	}
	return s.Data.AllKeys(), nil
}

// GetAllSettings merges all Settings and returns them as a map[string]interface{}.
func (s *Settings) GetAllSettings() (map[string]interface{}, error) {
	if s.Error != nil {
		return map[string]interface{}{}, fmt.Errorf("settings.GetAllSettings :: %s", s.Error)
	}
	return s.Data.AllSettings(), nil
}

// SubTree returns a new settings instance representing a sub tree of this instance.
// SubTree is case-insensitive for a key.
func (s *Settings) SubTree(prefix string) *Settings {
	s.Data = s.Data.Sub(prefix)
	return s
}

// GetSettingsFileNames returns the name of all settings files, whence settings manager was initialized.
func (s *Settings) GetSettingsFileNames() ([]string, error) {
	if s.Error != nil {
		return nil, fmt.Errorf("settings.GetSettingsFileNames :: %s", s.Error)
	}
	return s.fileNames, nil
}

// Reload once it's called, will re-read the settings data.
func (s *Settings) Reload() {
	s.mux.Lock()

	content := s.content
	s.Data = viper.New()

	if content != "" {
		s.Data = NewFromContent(content).Data
	}

	for _, fileName := range s.fileNames {
		s.Data = s.load(fileName).Data
	}

	s.mux.Unlock()
}

// AutoReload watching for settings file changes in the background
// and reloads configuration if needed.
func (s *Settings) AutoReload() {
	for _, fileName := range s.fileNames {
		s.Data.SetConfigFile(fileName)
		s.Data.WatchConfig()
		triggerReload(s)
	}
}
