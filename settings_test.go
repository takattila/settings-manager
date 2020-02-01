package settings

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type (
	unitConfSuite struct {
		suite.Suite
	}
)

func (u unitConfSuite) TestInitFromSource() {
	initTestOk()

	sm := NewFromContent(testYamlContent)
	u.Equal(nil, sm.Error)

	sm = NewFromContent(testJSONContent)
	u.Equal(nil, sm.Error)

	sm = NewFromContent(testBadYamlContent)
	u.Equal("settings.NewFromContent :: unsupported content type", fmt.Sprint(sm.Error))

	resetTest()
}

func (u unitConfSuite) TestInit() {
	initTestOk()

	sm := New(testYamlFilePAth)
	u.Equal(nil, sm.Error)

	sm = New(testDirPath)
	u.Equal(nil, sm.Error)

	resetTest()
}

func (u unitConfSuite) TestMerge() {
	initTestOk()

	err := ioutil.WriteFile(testJsonFileOtherPAth, []byte(testJSONContentOther), os.ModePerm)
	u.Equal(nil, err)

	sm := New(testYamlFilePAth).Merge(testJsonFileOtherPAth)
	u.Equal(nil, sm.Error)

	v, err := sm.Get("other.content.string")
	u.Equal(nil, err)
	u.Equal("text", v)

	v2, err := sm.Get("email.to")
	u.Equal(nil, err)
	u.Equal("user@gmail.com", v2)

	resetTest()
}

func (u unitConfSuite) TestGetAllKeys() {
	err := os.Mkdir(testDirPath, os.ModePerm)
	u.Equal(nil, err)

	err = ioutil.WriteFile(testYamlFileOtherPAth, []byte(testYamlContentOther), os.ModePerm)
	u.Equal(nil, err)

	sm := New(testYamlFileOtherPAth)
	u.Equal(nil, sm.Error)

	v, err := sm.GetAllKeys()
	u.Equal(nil, err)
	u.Contains(v, "other.content.int")
	u.Contains(v, "other.content.string")
	u.Contains(v, "other.content.bool")

	resetTest()

	err = os.Mkdir(testDirPath, os.ModePerm)
	u.Equal(nil, err)

	initTestBad()

	_, err = New(testBadYamlFilePAth).GetAllKeys()
	u.Equal("settings.GetAllKeys :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitConfSuite) TestGetAllSettings() {
	err := os.Mkdir(testDirPath, os.ModePerm)
	u.Equal(nil, err)

	err = ioutil.WriteFile(testJsonFileOtherPAth, []byte(testJSONContentOther), os.ModePerm)
	u.Equal(nil, err)

	sm := New(testJsonFileOtherPAth)
	u.Equal(nil, sm.Error)

	v, err := sm.GetAllSettings()
	u.Equal(nil, err)

	content := v["other"].(map[string]interface{})["content"].(map[string]interface{})

	u.Equal(true, content["bool"])
	u.Equal(float64(1), content["int"])
	u.Equal("text", content["string"])

	resetTest()

	err = os.Mkdir(testDirPath, os.ModePerm)
	u.Equal(nil, err)

	initTestBad()

	_, err = New(testBadYamlFilePAth).GetAllSettings()
	u.Equal("settings.GetAllSettings :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitConfSuite) TestAddPrefix() {
	initTestOk()

	s, err := New(testYamlFilePAth).SubTree("a.b").GetAllSettings()
	u.Equal(nil, err)
	u.Equal(map[string]interface{}{"c": []interface{}{1, 2, 0, 4}}, s)

	resetTest()
}

func (u unitConfSuite) TestGetSettingsFileNames() {
	initTestOk()

	fileNames, err := New(testYamlFilePAth).GetSettingsFileNames()
	u.Equal(nil, err)
	u.Equal([]string{"settings/test.yaml"}, fileNames)

	err = os.Chmod(testYamlFilePAth, os.ModeExclusive)
	u.Equal(nil, err)

	_, err = New(testYamlFilePAth).GetSettingsFileNames()
	u.Equal(`settings.GetSettingsFileNames :: open ./settings/test.yaml: permission denied`, fmt.Sprint(err))

	err = os.Chmod(testYamlFilePAth, os.ModePerm)
	u.Equal(nil, err)

	resetTest()
}

func (u unitConfSuite) TestReload() {
	initTest()

	sm := New(testYamlFilePAth)

	oldValue := `name: ExampleService`
	newValue := `name: NewApp`

	content := strings.ReplaceAll(testYamlContent, oldValue, newValue)
	saveFile(u, testYamlFilePAth, content)

	sm.Reload()

	v, err := sm.Get("service.name")
	u.Equal(nil, err)
	u.Equal("NewApp", v)

	resetTest()
}

func (u unitConfSuite) TestAutoReload() {
	initTest()

	err := ioutil.WriteFile(testYamlFileOtherPAth, []byte(testYamlContentOther), os.ModePerm)
	u.Equal(nil, err)

	wg := sync.WaitGroup{}
	wg.Add(1)

	tmpTriggerReload := triggerReload
	triggerReload = func(s *Settings) {
		go func() {
			for {
				s.Reload()
				time.Sleep(10 * time.Millisecond)

				wg.Done()
				return
			}
		}()
	}

	sm := NewFromContent(testYamlContent).Merge(testYamlFileOtherPAth)
	sm.AutoReload()

	oldValue := `int: 1`
	newValue := `int: 1000`

	content := strings.ReplaceAll(testYamlContentOther, oldValue, newValue)
	saveFile(u, testYamlFileOtherPAth, content)
	wg.Wait()

	v, err := sm.Get("other.content.int")
	u.Equal(nil, err)
	u.Equal(1000, v)

	v, err = sm.Get("service.name")
	u.Equal(nil, err)
	u.Equal("ExampleService", v)

	resetTest()

	triggerReload = tmpTriggerReload
}

func saveFile(u unitConfSuite, path, content string) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	u.Equal(nil, err)

	defer func() {
		_ = file.Close()
	}()

	// new writer w/ default 4096 buffer size
	w := bufio.NewWriter(file)

	_, err = w.WriteString(content + "\n")
	u.Equal(nil, err)

	err = w.Flush()
	u.Equal(nil, err)
}

func TestConfUnitSuite(t *testing.T) {
	suite.Run(t, new(unitConfSuite))
}
