package tinnitus

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"time"
)

var Config *TinnitusConfig = &TinnitusConfig{}

type TinnitusConfig struct {
	SC            SCConfig      `yaml:"sc"`
	RPC           RPCConfig     `yaml:"rpc"`
	Logger        zap.Config    `yaml:"logger"`
	SleepDuration time.Duration `yaml:"sleep_duration"`
}

type SCConfig struct {
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
}

type RPCConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func InitConfig() {
	var configPath string

	if Mode == TestingTinnitusMode {
		configPath = path.Join(PackagePath(), "config/config.testing.yaml")
	} else {
		configPath = *Flags.Config
	}

	configFile, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatalf("Unable to read the config file, got %s", err)
	}

	err = yaml.Unmarshal(configFile, Config)

	if err != nil {
		log.Fatalf("Unable to parse the config file, got %s", err)
	}
}
