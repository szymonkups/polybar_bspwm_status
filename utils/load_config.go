package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LabelState struct {
	Focused string
	Blurred string
	Inactive string
}

type Config struct {
	MonitorIndex int
	Labels struct {
		EmptyDesktop LabelState
		OccupiedDesktop LabelState
	}
}

func LoadConfig() (*Config, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	jsonFile, err := os.Open(dir + "/config.json")

	if err != nil {
		return nil, err
	}

	bytesRead, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	config := Config{}
	err = json.Unmarshal(bytesRead, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}