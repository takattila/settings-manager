package settings

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type (
	unitHelpersSuite struct {
		suite.Suite
	}
)

func (u unitHelpersSuite) TestTriggerReload() {
	initTestOk()

	sm := New(testYamlFilePAth)
	sm.Data.SetConfigFile(testYamlFilePAth)
	sm.Data.WatchConfig()

	saveFileHelper(u, testYamlFilePAth, testYamlContent)
	triggerReload(sm)

	time.Sleep(10 * time.Millisecond)

	v, err := sm.Get("service.name")
	u.Equal(nil, err)
	u.Equal("ExampleService", v)

	resetTest()
}

func (u unitHelpersSuite) TestLoad() {
	initTestOk()

	err := os.Chmod(testYamlFilePAth, os.ModeExclusive)
	u.Equal(nil, err)

	s := Settings{}
	s.Data = viper.New()

	sm := s.load(testFilePAth)
	u.Equal("open settings/test.yaml: permission denied", fmt.Sprint(sm.Error))

	resetTest()
}

func (u unitHelpersSuite) TestCheckErrors() {
	initTestOk()

	s := Settings{}
	s.Data = viper.New()

	sm := s.load(testFilePAth).checkErrors("key.not.exists", "checkErrors")
	u.Equal("settings.checkErrors :: key.not.exists :: cannot find value in configuration", fmt.Sprint(sm.Error))

	resetTest()
}

func (u unitHelpersSuite) TestCheckType() {
	initTestOk()

	s := Settings{}
	s.Data = viper.New()

	sm := s.load(testFilePAth).checkType("email.server.port", "checkType", reflect.String)
	u.Equal("settings.checkType :: the value of key: email.server.port :: should be type: string, not: int", fmt.Sprint(sm.Error))

	resetTest()
}

func (u unitHelpersSuite) TestCheck() {
	initTestOk()

	s := Settings{}
	s.Data = viper.New()

	sm := s.load(testFilePAth).check("email.server.port", "check", reflect.String)
	u.Equal("settings.check :: the value of key: email.server.port :: should be type: string, not: int", fmt.Sprint(sm.Error))

	resetTest()
}

func (u unitHelpersSuite) TestIsDirectory() {
	initTest()

	chk := isDirectory(testFilePAth)
	u.Equal(true, chk)

	chk = isDirectory(testYamlFilePAth)
	u.Equal(false, chk)

	chk = isDirectory("")
	u.Equal(false, chk)

	resetTest()
}

func (u unitHelpersSuite) TestValidateExtension() {
	ext := jsonExtension
	supported := ext.validateExtension()
	u.Equal(true, supported)

	ext = yamlExtensionLong
	supported = ext.validateExtension()
	u.Equal(true, supported)

	ext = yamlExtensionShort
	supported = ext.validateExtension()
	u.Equal(true, supported)

	ext = ".ini"
	supported = ext.validateExtension()
	u.Equal(false, supported)
}

func TestHelperUnitSuite(t *testing.T) {
	suite.Run(t, new(unitHelpersSuite))
}

func saveFileHelper(u unitHelpersSuite, path, content string) {
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

func initTest() {
	resetTest()
	initTestOk()
	initTestBad()
}

func initTestOk() {
	_ = os.Mkdir(testDirPath, os.ModePerm)

	err := ioutil.WriteFile(testYamlFilePAth, []byte(testYamlContent), os.ModePerm)
	checkErr("initTest", err)
}

func initTestBad() {
	_ = os.Mkdir(testDirPath, os.ModePerm)

	err := ioutil.WriteFile(testBadYamlFilePAth, []byte(testBadYamlContent), os.ModePerm)
	checkErr("initTest", err)
}

func resetTest() {
	err := os.RemoveAll(testDirPath)
	checkErr("initTest", err)
}

func checkErr(funcName string, err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	testYamlContent = `
service:
  name: ExampleService
server:
  timeout:
    read: 5
    write: 10
    idle: 10
email:
  to: user@gmail.com
  server:
    user: user.name@gmail.com
    address: smtp.gmail.com
    port: 587
intSlice:
  - 0
  - 1
  - 2
  - 3
stringSlice:
  - 0
  - 1
  - other value
  - 3
float64: 123131232132113211564564456
dotted.key:
  value: it is a value of a dotted key
a:
  b:
    c:
      - 1
      - 2
      - 0
      - 4
restart:
  always: true
environment:
  dev:
    server:
      url: 'http://example-server.com:7878'
    app1:
      level: debug
      host: example-service.com
      port: 12201
      network: udp
      compressed: false
      packet_size: 262144
      enabled: true
    app2:
      host: 0.0.0.0
      port: 6125
      prefix: ExampleService
  prod:
    server:
      url: 'http://example-server.com:8888'
    app1:
      level: info
      host: example-service.com
      port: 12201
      network: udp
      compressed: true
      packet_size: 262144
      enabled: false
    app2:
      host: 0.0.0.0
      port: 8125
      prefix: ExampleService
  test:
    server:
      url: 'http://example-server.com:7878'
    app1:
      level: info
      host: example-service.com
      port: 12201
      network: udp
      compressed: true
      packet_size: 262144
      enabled: false
    app2:
      host: 0.0.0.0
      port: 7125
      prefix: ExampleService
`

	testYamlContentOther = `
other:
  content:
    int: 1
    string: 'text'
    bool: true
`

	testJSONContent = `
{
  "service": {
    "name": "ExampleService"
  },
  "server": {
    "timeout": {
      "read": 5,
      "write": 10,
      "idle": 10
    }
  },
  "email": {
    "to": "user@gmail.com",
    "server": {
      "user": "user.name@gmail.com",
      "address": "smtp.gmail.com",
      "port": 587
    }
  },
  "intSlice": [
    0,
    1,
    2,
    3
  ],
  "stringSlice": [
    0,
    1,
    "other value",
    3
  ],
  "float64": 1.231312321321132e+26,
  "dotted.key": {
    "value": "it is a value of a dotted key"
  },
  "a": {
    "b": {
      "c": [
        1,
        2,
        0,
        4
      ]
    }
  },
  "restart": {
    "always": true
  },
  "environment": {
    "dev": {
      "server": {
        "url": "http://example-server.com:7878"
      },
      "app1": {
        "level": "debug",
        "host": "example-service.com",
        "port": 12201,
        "network": "udp",
        "compressed": false,
        "packet_size": 262144,
        "enabled": true
      },
      "app2": {
        "host": "0.0.0.0",
        "port": 6125,
        "prefix": "ExampleService"
      }
    },
    "prod": {
      "server": {
        "url": "http://example-server.com:8888"
      },
      "app1": {
        "level": "info",
        "host": "example-service.com",
        "port": 12201,
        "network": "udp",
        "compressed": true,
        "packet_size": 262144,
        "enabled": false
      },
      "app2": {
        "host": "0.0.0.0",
        "port": 8125,
        "prefix": "ExampleService"
      }
    },
    "test": {
      "server": {
        "url": "http://example-server.com:7878"
      },
      "app1": {
        "level": "info",
        "host": "example-service.com",
        "port": 12201,
        "network": "udp",
        "compressed": true,
        "packet_size": 262144,
        "enabled": false
      },
      "app2": {
        "host": "0.0.0.0",
        "port": 7125,
        "prefix": "ExampleService"
      }
    }
  }
}
`

	testJSONContentOther = `
{
 "other": {
   "content": {
     "int": 1,
     "string": "text",
     "bool": true
   }
 }
}
`

	testBadYamlContent = `/* BAD YAML * \`

	testBadYamlFilePAth = "./settings/bad.yaml"

	testFilePAth = "./settings"

	testYamlFilePAth = "./settings/test.yaml"

	testJsonFilePAth = "./settings/test.json"

	testJsonFileOtherPAth = "./settings/other.json"

	testYamlFileOtherPAth = "./settings/other.yaml"

	testDirPath = "./settings"
)
