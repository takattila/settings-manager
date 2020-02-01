package settings_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/takattila/settings-manager"
)

func ExampleSettings_Get() {
	file := "example_app1.yaml"
	content := `{ "app": { "string": "value", "int": 1 } }`

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.Get("app.string")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)

	v, err = sm.Get("app.int")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)

	// Output:
	// value
	// 1
}

func ExampleSettings_GetBool() {
	file := "example_app1.yaml"
	content := "app:\n  bool:  true"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetBool("app.bool")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("value: %t, type: %T\n", v, v)

	// Output: value: true, type: bool
}

func ExampleSettings_GetFloat64() {
	file := "example_app1.yaml"
	content := "app:\n  float64:  123131232132113211564564456"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetFloat64("app.float64")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("value: %b, type: %T\n", v, v)

	// Output: value: 7167181007803488p+34, type: float64
}

func ExampleSettings_GetInt() {
	file := "example_app1.yaml"
	content := "app:\n  int:  100"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetInt("app.int")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("value: %d, type: %T\n", v, v)

	// Output: value: 100, type: int
}

func ExampleSettings_GetIntSlice() {
	file := "example_app1.yaml"
	content := `{ "app": { "int_slice": [ 1, 2 ] } }`

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetIntSlice("app.int_slice")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("value: %d, type: %T\n", v, v)

	// Output: value: [1 2], type: []int
}

func ExampleSettings_GetString() {
	file := "example_app1.yaml"
	content := "app:\n  string: some text"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetString("app.string")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("value: %s, type: %T\n", v, v)

	// Output: value: some text, type: string
}

func ExampleSettings_GetStringMap() {
	file := "example_app1.yaml"
	content := `{ "app": { "string_map": { "one": "value1", "two": "value2" } } }`

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetStringMap("app.string_map")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("type: %T\n", v)
	fmt.Printf("app.string_map.one: %s\n", v["one"])
	fmt.Printf("app.string_map.two: %s\n", v["two"])

	// Output:
	// type: map[string]interface {}
	// app.string_map.one: value1
	// app.string_map.two: value2
}

func ExampleSettings_GetStringMapString() {
	file := "example_app1.yaml"
	content := `{ "app": { "string_map": { "one": "value1", "two": "value2" } } }`

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetStringMapString("app.string_map")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("type: %T\n", v)
	fmt.Printf("app.string_map.one: %s\n", v["one"])
	fmt.Printf("app.string_map.two: %s\n", v["two"])

	// Output:
	// type: map[string]string
	// app.string_map.one: value1
	// app.string_map.two: value2
}

func ExampleSettings_GetStringSlice() {
	file := "example_app1.yaml"
	content := `{ "app": { "string_slice": [ "value1", "value2" ] } }`

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetStringSlice("app.string_slice")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("app.string_slice: %s, type: %T\n", v, v)

	// Output:
	// app.string_slice: [value1 value2], type: []string
}

func ExampleSettings_GetTime() {
	file := "example_app1.yaml"
	content := "time:\n  seconds: 10"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetTime("time.seconds")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("time.seconds: %d, type: %T\n", v.Second(), v)

	// Output:
	// time.seconds: 10, type: time.Time
}

func ExampleSettings_GetDuration() {
	file := "example_app1.yaml"
	content := "time:\n  duration: 10"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	v, err := sm.GetDuration("time.duration")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("time.duration: %d, type: %T\n", v, v)

	// Output:
	// time.duration: 10, type: time.Duration
}

func ExampleSettings_IsSet() {
	file := "example_app1.yaml"
	content := "app1:\n  key: value"

	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sm := settings.New(file)

	isSet, err := sm.IsSet("app1.key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("app1.key: %t, type: %T\n", isSet, isSet)

	isSet, err = sm.IsSet("app2.key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("app2.key: %t, type: %T\n", isSet, isSet)

	// Output:
	// app1.key: true, type: bool
	// app2.key: false, type: bool
}
