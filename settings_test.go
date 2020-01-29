package settings

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type (
	unitConfSuite struct {
		suite.Suite
	}
)

func (u unitConfSuite) TestInitFromSource() {
	initTestOk()

	sm := NewFromSource(testYamlContent)
	u.Equal(nil, sm.Error)

	sm = NewFromSource(testJSONContent)
	u.Equal(nil, sm.Error)

	sm = NewFromSource(testBadYamlContent)
	u.Equal("settings.NewFromSource :: unsupported content type", fmt.Sprint(sm.Error))

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

	err := ioutil.WriteFile(testJsonlFileOtherPAth, []byte(testJSONContentOther), os.ModePerm)
	u.Equal(nil, err)

	sm := New(testYamlFilePAth).Merge(testJsonlFileOtherPAth)
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

	err = ioutil.WriteFile(testYamllFileOtherPAth, []byte(testYamlContentOther), os.ModePerm)
	u.Equal(nil, err)

	sm := New(testYamllFileOtherPAth)
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

	err = ioutil.WriteFile(testJsonlFileOtherPAth, []byte(testJSONContentOther), os.ModePerm)
	u.Equal(nil, err)

	sm := New(testJsonlFileOtherPAth)
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

	s, err := New(testYamlFilePAth).AddPrefix("a.b").GetAllSettings()
	u.Equal(nil, err)
	u.Equal(map[string]interface{}{"c": []interface{}{1, 2, 0, 4}}, s)

	resetTest()
}

func TestConfUnitSuite(t *testing.T) {
	suite.Run(t, new(unitConfSuite))
}
