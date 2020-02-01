# GO - Settings Manager

[![Test](https://github.com/takattila/settings-manager/workflows/Test/badge.svg?branch=master)](https://github.com/takattila/settings-manager/actions?query=workflow:Test)
[![Coverage Status](https://coveralls.io/repos/github/takattila/settings-manager/badge.svg?branch=master)](https://coveralls.io/github/takattila/settings-manager?branch=master)
[![GOdoc](https://img.shields.io/badge/godoc-reference-orange)](https://godoc.org/github.com/takattila/settings-manager)
[![Version](https://img.shields.io/badge/dynamic/json.svg?label=version&url=https://api.github.com/repos/takattila/settings-manager/releases/latest&query=tag_name)](https://github.com/takattila/settings-manager/releases)

This package was made, to easily get needed settings from a file.

- Supported file types are: **'json'**, and **'yaml'**.
- The configuration keys are **case insensitive**.

This package uses [github.com/spf13/viper](https://github.com/spf13/viper)

## Table of Contents

* [Example usage](#example-usage)
   * [Initialization](#initialization)
   * [Merge configuration with other settings file](#merge-configuration-with-other-settings-file)
   * [Initialize settings from a given content](#initialize-settings-from-a-given-content)
   * [Get all keys from the settings](#get-all-keys-from-the-settings)
   * [Get all settings](#get-all-settings)
   * [Add a sub tree](#add-a-sub-tree)
   * [Type assertions](#type-assertions)
   * [Reload the settings data manually](#reload-the-settings-data-manually)
   * [Automatic reload the settings data in the background](#automatic-reload-the-settings-data-in-the-background)

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
sm := settings.NewFromContent(content)
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
sm := settings.NewFromContent(content)
keys, err := sm.GetAllKeys()
if err != nil {
    log.Fatal(err)
}

// Output:
// 2020/01/29 19:09:52 keys [other content.int content.string content.bool]
log.Println("keys", keys)
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
sm := settings.NewFromContent(content)
allSettings, err := sm.GetAllSettings()
if err != nil {
    log.Fatal(err)
}

// Output:
// 2020/01/29 19:11:58 allSettings map[content:map[bool:true int:1 string:text]]
log.Println("allSettings", allSettings)
```

[Back to top](#table-of-contents)

### Add a sub tree

Returning a new settings instance representing a sub tree of this instance.
SubTree is case-insensitive for a key.

```go
var content = `
a:
  b:
    c:
      value: 1
`

sm := settings.NewFromContent(content)

sm = sm.SubTree("a.b")

intValue, err := sm.GetInt("c.value")
if err != nil {
    log.Fatal(err)
}

// Output:
// 2020/02/02 15:52:55 allSettings 1
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

### Reload the settings data manually

Re-read the settings data by calling the Reload function.

```go
content := `
config:
  config_key:  config_value`

err := ioutil.WriteFile("./example/settings/example_config.yaml", []byte(content), os.ModePerm)
if err != nil {
    log.Fatal(err)
}

sm := settings.New("./example/settings/example_config.yaml")

v, err := sm.Get("config.config_key")
if err != nil {
    log.Fatal(err)
}

// Output:
// 2020/02/02 15:31:58 config_value
log.Println(v)

content = strings.ReplaceAll(content, "config_key:  config_value", "foo:  bar")
err = ioutil.WriteFile("./example/settings/example_config.yaml", []byte(content), os.ModePerm)
if err != nil {
    log.Fatal(err)
}

// Reload the configuration ...
sm.Reload()

v, err = sm.Get("config.foo")
if err != nil {
    log.Fatal(err)
}

// Output:
// 2020/02/02 15:31:58 bar
log.Println(v)
```

[Back to top](#table-of-contents)

### Automatic reload the settings data in the background

AutoReload is watching for settings file changes in the background and reloads configuration if needed.

```go
content := `
config:
  config_key:  config_value`

err := ioutil.WriteFile("./example/settings/example_config.yaml", []byte(content), os.ModePerm)
if err != nil {
    log.Fatal(err)
}

sm := settings.New("./example/settings/example_config.yaml")

// Activate the automatic reload function ...
sm.AutoReload()

v, err := sm.Get("config.config_key")
if err != nil {
    log.Fatal(err)
}

// Output:
// 2020/02/02 15:42:49 config_value
log.Println(v)

content = strings.ReplaceAll(content, "config_key:  config_value", "foo:  bar")
err = ioutil.WriteFile("./example/settings/example_config.yaml", []byte(content), os.ModePerm)
if err != nil {
    log.Fatal(err)
}

time.Sleep(5 * time.Millisecond)

v, err = sm.Get("config.foo")
if err != nil {
    log.Fatal(err)
}

// Output:
// 2020/02/02 15:42:49 settings.AutoReload settings reloaded
// 2020/02/02 15:42:49 bar
log.Println(v)

content = strings.ReplaceAll(content, "foo:  bar", "config_key:  config_value")
err = ioutil.WriteFile("./example/settings/example_config.yaml", []byte(content), os.ModePerm)
if err != nil {
    log.Fatal(err)
}

time.Sleep(5 * time.Millisecond)

v, err = sm.Get("config.config_key")
if err != nil {
    log.Fatal(err)
}

// Output:
// 2020/02/02 15:42:49 settings.AutoReload settings reloaded
// 2020/02/02 15:42:49 config_value
log.Println(v)
```

[Back to top](#table-of-contents)
