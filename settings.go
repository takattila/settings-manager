// This package was made, to easily get needed settings from a file.
// Supported file types are: **'*.json'** and **'*.yaml'**.
//
// This package uses https://github.com/spf13/viper: Copyright © 2014 Steve Francia <spf@spf13.com>.
package settings

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

type Settings struct {
	Data  *viper.Viper
	Error error
}

// New initializes settings from a file or from multiple files under given directory.
func New(settingsFile string) *Settings {
	s := &Settings{}
	s.Data = viper.New()
	return s.load(settingsFile)
}

// NewFromSource initializes settings from a given content.
func NewFromSource(content string) *Settings {
	s := &Settings{}
	s.Data = viper.New()
	ext := getExtensionByContent(content)
	if ext == "unsupported" {
		return &Settings{Error: fmt.Errorf("settings.NewFromSource :: unsupported content type")}
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

// AddPrefix returns a new settings instance representing a sub tree of this instance.
// AddPrefix is case-insensitive for a key.
func (s *Settings) AddPrefix(prefix string) *Settings {
	s.Data = s.Data.Sub(prefix)
	return s
}
