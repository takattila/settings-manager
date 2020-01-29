package settings

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type (
	unitCastSuite struct {
		suite.Suite
	}
)

func (u unitCastSuite) TestGet() {
	initTest()

	v, err := New(testYamlFilePAth).Get("service.name")
	u.Equal(nil, err)
	u.Equal("ExampleService", v)

	_, err = New(testBadYamlFilePAth).Get("service.name")
	u.Equal("settings.Get :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetBool() {
	initTest()

	v, err := New(testYamlFilePAth).GetBool("restart.always")
	u.Equal(nil, err)
	u.Equal(true, v)

	_, err = New(testBadYamlFilePAth).GetBool("restart.always")
	u.Equal("settings.GetBool :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetFloat64() {
	initTest()

	v, err := New(testYamlFilePAth).GetFloat64("float64")
	u.Equal(nil, err)
	u.Equal(1.231312321321132e+26, v)

	_, err = New(testBadYamlFilePAth).GetFloat64("float64")
	u.Equal("settings.GetFloat64 :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetInt() {
	initTest()

	v, err := New(testYamlFilePAth).GetInt("server.timeout.read")
	u.Equal(nil, err)
	u.Equal(5, v)

	_, err = New(testBadYamlFilePAth).GetInt("server.timeout.read")
	u.Equal("settings.GetInt :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetIntSlice() {
	initTest()

	sm := New(testYamlFilePAth)

	v, err := sm.GetIntSlice("intSlice")
	u.Equal(nil, err)
	u.Equal([]int{0, 1, 2, 3}, v)

	_, err = sm.GetIntSlice("stringSlice")
	u.Equal("settings.GetIntSlice :: the value of key: stringSlice :: should be type: []int, not: []interface {}", fmt.Sprint(err))

	_, err = New(testBadYamlFilePAth).GetIntSlice("intSlice")
	u.Equal("settings.GetIntSlice :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetString() {
	initTest()

	v, err := New(testYamlFilePAth).GetString("email.server.user")
	u.Equal(nil, err)
	u.Equal("user.name@gmail.com", v)

	_, err = New(testBadYamlFilePAth).GetString("email.server.user")
	u.Equal("settings.GetString :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetStringMap() {
	initTest()

	v, err := New(testYamlFilePAth).GetStringMap("environment.prod.app1")
	u.Equal(nil, err)
	u.Equal(false, v["enabled"])

	u.Equal(true, v["compressed"])
	u.Equal("example-service.com", v["host"])
	u.Equal("info", v["level"])
	u.Equal(262144, v["packet_size"])
	u.Equal("udp", v["network"])
	u.Equal(12201, v["port"])

	_, err = New(testBadYamlFilePAth).GetStringMap("email.server.user")
	u.Equal("settings.GetStringMap :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetStringMapString() {
	initTest()

	v, err := New(testYamlFilePAth).GetStringMapString("environment.prod.app1")
	u.Equal(nil, err)
	u.Equal("false", v["enabled"])

	u.Equal("true", v["compressed"])
	u.Equal("example-service.com", v["host"])
	u.Equal("info", v["level"])
	u.Equal("262144", v["packet_size"])
	u.Equal("udp", v["network"])
	u.Equal("12201", v["port"])

	_, err = New(testBadYamlFilePAth).GetStringMapString("email.server.user")
	u.Equal("settings.GetStringMapString :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetStringSlice() {
	initTest()

	v, err := New(testYamlFilePAth).GetStringSlice("stringSlice")
	u.Equal(nil, err)
	u.Equal([]string{"0", "1", "other value", "3"}, v)

	_, err = New(testBadYamlFilePAth).GetStringSlice("stringSlice")
	u.Equal("settings.GetStringSlice :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetTime() {
	initTest()

	v, err := New(testYamlFilePAth).GetTime("server.timeout.write")
	u.Equal(nil, err)
	u.Equal(10, v.Second())

	_, err = New(testBadYamlFilePAth).GetTime("server.timeout.write")
	u.Equal("settings.GetTime :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestGetDuration() {
	initTest()

	v, err := New(testYamlFilePAth).GetDuration("server.timeout.write")
	u.Equal(nil, err)
	u.Equal(time.Duration(10), v)

	_, err = New(testBadYamlFilePAth).GetDuration("server.timeout.write")
	u.Equal("settings.GetDuration :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))

	resetTest()
}

func (u unitCastSuite) TestIsSet() {
	initTest()

	v, err := New(testYamlFilePAth).IsSet("environment.dev.app1.level")
	u.Equal(nil, err)
	u.Equal(true, v)

	v, err = New(testYamlFilePAth).IsSet("not.existent.key")
	u.Equal(nil, err)
	u.Equal(false, v)

	v, err = New(testBadYamlFilePAth).IsSet("server.timeout.write")
	u.Equal("settings.IsSet :: While parsing config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `/* BAD ...` into map[string]interface {}", fmt.Sprint(err))
	u.Equal(false, v)

	resetTest()
}

func TestCastUnitSuite(t *testing.T) {
	suite.Run(t, new(unitCastSuite))
}
