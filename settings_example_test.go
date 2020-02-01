package settings_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/takattila/settings-manager"
)

func ExampleNew() {
	file := "example_config.yaml"
	content := "config:\n  key:  value"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	AllSettings, err := sm.GetAllSettings()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(AllSettings)

	// Output: map[config:map[key:value]]
}

func ExampleNewFromContent() {
	content := "config:\n  key:  value"

	sm := settings.NewFromContent(content)

	AllSettings, err := sm.GetAllSettings()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(AllSettings)

	// Output: map[config:map[key:value]]
}

func ExampleSettings_Merge() {
	file1 := "example_app1.yaml"
	content1 := "app1:\n  key:  value1"

	err := ioutil.WriteFile(file1, []byte(content1), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file2 := "example_app2.yaml"
	content2 := "app2:\n  key:  value2"

	err = ioutil.WriteFile(file2, []byte(content2), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file1).Merge(file2)

	k, err := sm.Get("app1.key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k)

	k, err = sm.Get("app2.key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k)

	// Output:
	// value1
	// value2
}

func ExampleSettings_GetAllKeys() {
	content := "config:\n  key:  value"

	sm := settings.NewFromContent(content)
	keys, err := sm.GetAllKeys()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(keys)

	// Output: [config.key]
}

func ExampleSettings_GetAllSettings() {
	content := "config:\n  key:  value"

	sm := settings.NewFromContent(content)
	AllSettings, err := sm.GetAllSettings()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(AllSettings)

	// Output: map[config:map[key:value]]
}

func ExampleSettings_SubTree() {
	content := `{ "config": { "sub": { "tree": "value" } } }`

	sm := settings.NewFromContent(content)
	sm = sm.SubTree("config.sub")

	v, err := sm.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)

	// Output: value
}

func ExampleSettings_GetSettingsFileNames() {
	file1 := "example_app1.yaml"
	content1 := "app1:\n  key:  value1"

	err := ioutil.WriteFile(file1, []byte(content1), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file2 := "example_app2.yaml"
	content2 := "app2:\n  key:  value2"

	err = ioutil.WriteFile(file2, []byte(content2), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file1).Merge(file2)

	files, err := sm.GetSettingsFileNames()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(files)

	// Output: [example_app1.yaml example_app2.yaml]
}

func ExampleSettings_Reload() {
	file := "example_config.yaml"
	content := "config:\n  key:  value"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.Get("config.key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)

	content = strings.ReplaceAll(content, "key:  value", "foo:  bar")
	err = ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Reload the configuration ...
	sm.Reload()

	v, err = sm.Get("config.foo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)

	// Output:
	// value
	// bar
}

func ExampleSettings_AutoReload() {
	file := "example_config.yaml"
	content := "config:\n  key:  value"

	// Save content
	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	// Activate the automatic reload function ...
	sm.AutoReload()

	v, err := sm.Get("config.key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)

	// First update of the content
	content = strings.ReplaceAll(content, "key:  value", "foo:  bar")
	err = ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Millisecond)

	v, err = sm.Get("config.foo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)

	// Second update of the content
	content = strings.ReplaceAll(content, "foo:  bar", "key:  value")
	err = ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Millisecond)

	v, err = sm.Get("config.key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)

	// Output:
	// value
	// bar
	// value
}
