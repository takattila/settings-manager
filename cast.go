package settings

import (
	"fmt"
	"reflect"
	"time"
)

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Settings will checkErrors in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func (s *Settings) Get(key string) (interface{}, error) {
	if err := s.check(key, "Get").Error; err != nil {
		return nil, err
	}
	return s.Data.Get(key), nil
}

// GetBool returns the value associated with the key as a boolean.
func (s *Settings) GetBool(key string) (bool, error) {
	if err := s.check(key, "GetBool", reflect.Bool).Error; err != nil {
		return false, err
	}
	return s.Data.GetBool(key), nil
}

// GetFloat64 returns the value associated with the key as a float64.
func (s *Settings) GetFloat64(key string) (float64, error) {
	if err := s.check(key, "GetFloat64", reflect.Float64).Error; err != nil {
		return 0, err
	}
	return s.Data.GetFloat64(key), nil
}

// GetInt returns the value associated with the key as an integer.
func (s *Settings) GetInt(key string) (int, error) {
	if err := s.check(key, "GetInt", reflect.Int).Error; err != nil {
		return 0, err
	}
	return s.Data.GetInt(key), nil
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func (s *Settings) GetIntSlice(key string) ([]int, error) {
	if err := s.checkIntSlice(key); err != nil {
		return []int{}, err
	}
	return s.Data.GetIntSlice(key), nil
}

// GetString returns the value associated with the key as a string.
func (s *Settings) GetString(key string) (string, error) {
	if err := s.check(key, "GetString", reflect.String).Error; err != nil {
		return "", err
	}
	return s.Data.GetString(key), nil
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (s *Settings) GetStringMap(key string) (map[string]interface{}, error) {
	if err := s.check(key, "GetStringMap", reflect.Map).Error; err != nil {
		return map[string]interface{}{}, err
	}
	return s.Data.GetStringMap(key), nil
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (s *Settings) GetStringMapString(key string) (map[string]string, error) {
	if err := s.check(key, "GetStringMapString", reflect.Map).Error; err != nil {
		return map[string]string{}, err
	}
	return s.Data.GetStringMapString(key), nil
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (s *Settings) GetStringSlice(key string) ([]string, error) {
	if err := s.check(key, "GetStringSlice", reflect.Slice).Error; err != nil {
		return []string{}, err
	}
	return s.Data.GetStringSlice(key), nil
}

// GetTime returns the value associated with the key as time.
func (s *Settings) GetTime(key string) (time.Time, error) {
	if err := s.check(key, "GetTime", reflect.Int).Error; err != nil {
		return time.Time{}, err
	}
	return s.Data.GetTime(key), nil
}

// GetDuration returns the value associated with the key as a duration.
func (s *Settings) GetDuration(key string) (time.Duration, error) {
	if err := s.check(key, "GetDuration", reflect.Int).Error; err != nil {
		return 0, err
	}
	return s.Data.GetDuration(key), nil
}

// IsSet checks to see if the key has been set in any of the Data locations.
// IsSet is case-insensitive for a key.
func (s *Settings) IsSet(key string) (bool, error) {
	if s.Error != nil {
		return false, fmt.Errorf("settings.IsSet :: %s", s.Error)
	}
	return s.Data.IsSet(key), nil
}
