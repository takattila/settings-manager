package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"

	"gopkg.in/yaml.v2"
)

type supportedExtension string

const (
	jsonExtension      supportedExtension = ".json"
	yamlExtensionLong  supportedExtension = ".yaml"
	yamlExtensionShort supportedExtension = ".yml"
)

var triggerReload = func(s *Settings) {
	s.Data.OnConfigChange(func(in fsnotify.Event) {
		s.Reload()
		log.Println("settings.AutoReload", "settings reloaded")
	})
}

func (s *Settings) load(settingsFile string) *Settings {
	if isDirectory(settingsFile) {
		for _, file := range listFilesUnderDirectory(settingsFile) {
			if err := s.load(file).Error; err != nil {
				return &Settings{Error: err}
			}
		}
	} else {
		b, err := ioutil.ReadFile(settingsFile)
		if err != nil {
			return &Settings{Error: err}
		}
		s.Data.SetConfigType(getExtensionByFileName(settingsFile))

		err = s.Data.MergeConfig(bytes.NewBuffer(b))
		if err != nil {
			return &Settings{Error: err}
		}
		s.appendFileName(settingsFile)
	}
	return s
}

func (s *Settings) checkErrors(key, funcName string) *Settings {
	if s.Error != nil {
		return &Settings{Error: fmt.Errorf("settings.%s :: %s", funcName, s.Error)}
	} else if !s.Data.IsSet(key) {
		return &Settings{Error: fmt.Errorf("settings.%s :: %s :: cannot find value in configuration", funcName, key)}
	}
	return s
}

func (s *Settings) checkType(key, funcName string, kind reflect.Kind) *Settings {
	if reflect.TypeOf(s.Data.Get(key)).Kind() != kind {
		return &Settings{Error: fmt.Errorf("settings.%s :: the value of key: %s :: should be type: %s, not: %s",
			funcName,
			key,
			kind,
			reflect.TypeOf(s.Data.Get(key)))}
	}
	return s
}

func (s *Settings) check(key, funcName string, kind ...reflect.Kind) *Settings {
	if err := s.checkErrors(key, funcName).Error; err != nil {
		return &Settings{Error: err}
	} else if len(kind) > 0 {
		if err := s.checkType(key, funcName, kind[0]).Error; err != nil {
			return &Settings{Error: err}
		}
	}
	return s
}

func (s *Settings) checkIntSlice(key string) error {
	funcName := "GetIntSlice"

	if err := s.check(key, funcName, reflect.Slice).Error; err != nil {
		return err
	}

	t := s.Data.Get(key)

	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		sl := reflect.ValueOf(t)

		for i := 0; i < sl.Len(); i++ {
			_, ok := sl.Index(i).Interface().(int)
			if !ok {
				return fmt.Errorf("settings.%s :: the value of key: %s :: should be type: %s, not: %s",
					funcName,
					key,
					"[]int",
					reflect.TypeOf(s.Data.Get(key)))
			}
		}
	}
	return nil
}

func listFilesUnderDirectory(dir string) (files []string) {
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !isDirectory(path) {
			if supportedExtension(filepath.Ext(path)).validateExtension() {
				files = append(files, path)
			}
		}
		return nil
	})
	return
}

func getExtensionByContent(source string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(source), &obj); err == nil {
		return "json"
	}
	if err := yaml.Unmarshal([]byte(source), &obj); err == nil {
		return "yaml"
	}
	return "unsupported"
}

func getExtensionByFileName(fileName string) string {
	return strings.Replace(filepath.Ext(fileName), ".", "", -1)
}

func isDirectory(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := fi.Mode()
	if mode.IsDir() {
		return true
	}
	return false
}

func (s *Settings) appendFileName(fileName string) {
	s.fileNames = append(s.fileNames, filepath.Clean(fileName))
	s.fileNames = makeUniqueSlice(s.fileNames)
}

func makeUniqueSlice(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}
	return us
}

func (e supportedExtension) validateExtension() bool {
	switch e {
	case jsonExtension, yamlExtensionLong, yamlExtensionShort:
		return true
	}
	return false
}
