# GO - Settings Manager

[![Test](https://github.com/takattila/settings-manager/workflows/Test/badge.svg?branch=master)](https://github.com/takattila/settings-manager/actions?query=workflow:Test)
[![Coverage Status](https://coveralls.io/repos/github/takattila/settings-manager/badge.svg?branch=master)](https://coveralls.io/github/takattila/settings-manager?branch=master)
[![GOdoc](https://img.shields.io/badge/godoc-reference-orange)](https://godoc.org/github.com/takattila/settings-manager)
[![Version](https://img.shields.io/badge/dynamic/json.svg?label=version&url=https://api.github.com/repos/takattila/settings-manager/releases/latest&query=tag_name)](https://github.com/takattila/settings-manager/releases)

This package was made, to easily get needed settings from a file.

- Supported file types are: **'json'** and **'yaml'**.
- The configuration keys are **case insensitive**.

This package uses [github.com/spf13/viper](https://github.com/spf13/viper)

## Table of Contents

* [Example usage](#example-usage)
   * [Initialization](#initialization)
   * [Merge configuration with other settings file](#merge-configuration-with-other-settings-file)
   * [Initialize settings from a given content](#initialize-settings-from-a-given-content)
   * [Get all keys from the settings](#get-all-keys-from-the-settings)
   * [Get all settings](#get-all-settings)
   * [Add prefix](#add-prefix)
   * [Type assertions](#type-assertions)

## Example usage

### Initialization

Initialize settings from a file or from multiple files under given directory.

```go
sm := settings.New("./example/settings/config.yaml")

// ... or by passing directory path:

sm := settings.New("./example/settings")
```

[Back to top](#table-of-contents)

### Merge configuration with other settings file

Merge initialized settings with a given file or directory.

```go
sm := settings.New("./example/settings/config.yaml").
	Merge("./example/settings/test.yaml")
```

[Back to top](#table-of-contents)

### Initialize settings from a given content

Initialize settings from a given content.

```go
content := `
other:
  content:
    int: 1
    string: 'text'
    bool: true
`
sm := settings.New(NewFromSource)
```

[Back to top](#table-of-contents)

### Get all keys from the settings

Return all keys holding a value, regardless of where they are set.

```go
content := `
other:
  content:
    int: 1
    string: 'text'
    bool: true
`
sm := settings.New(NewFromSource)
keys, err := sm.GetAllKeys()
if err != nil {
    log.Fatal(err)
}

log.Println("keys", keys)

// Output:
// 2020/01/29 19:09:52 keys [other content.int content.string content.bool]
```

[Back to top](#table-of-contents)

### Get all settings

Merge all Settings and return them.

```go
content := `
other:
  content:
    int: 1
    string: 'text'
    bool: true
`
sm := settings.New(NewFromSource)
allSettings, err := sm.GetAllSettings()
if err != nil {
    log.Fatal(err)
}

log.Println("allSettings", allSettings)

// Output:
// 2020/01/29 19:11:58 allSettings map[content:map[bool:true int:1 string:text]]
```

[Back to top](#table-of-contents)

### Add prefix

Returning a new settings instance representing a sub tree of this instance.
AddPrefix is case-insensitive for a key.

```go
var content = `
a:
  b:
    c:
      value: 1
`

sm := settings.NewFromSource(content)

sm = sm.AddPrefix("a.b")

intValue, err := sm.GetInt("c.value")
if err != nil {
    log.Fatal(err)
}

log.Println("allSettings", intValue)
```

[Back to top](#table-of-contents)

### Type assertions

```go
sm := settings.New("./example/settings/settings.yaml")

interfaceKey, err := sm.Get("interface.key")
boolKey, err := sm.GetBool("bool.key")
float64Key, err := sm.GetFloat64("float64.key")
intKey, err := sm.GetInt("int.key")
stringKey, err := sm.GetString("string.key")
intSliceKey, err := sm.GetIntSlice("int.slice.key")
stringMapKey, err := sm.GetStringMap("string.map.key")
stringMapKey, err := sm.GetStringMapString("string.map.key")
stringSliceKey, err := sm.GetStringSlice("string.slice.key")
timeKey, err := sm.GetTime("time.key")
timeDurationKey, err := sm.GetDuration("time.duration.key")
```

[Back to top](#table-of-contents)
