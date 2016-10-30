package skyline

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

var Debug bool

type FilterItemConf struct {
	ItemNamePrefix string  `json:"item_name_prefix"`
	Cycle          int64   `json:"cycle"`
	MatchStr       string  `json:"match_str"`
	Threshold      float64 `json:"threshold"`
}

type MonitorConf struct {
	LogNamePrefix string           `json:"log_name_prefix"`
	LogFilePath   string           `json:"log_file_path"`
	FilterItems   []FilterItemConf `json:"filter_items"`
}

type WarningConf struct {
	WarningName   string `json:"warning_name"`
	Formula       string `json:"formula"`
	WarningFilter string `json:"warning_filter"`
	AlertName     string `json:"alert_name"`
	AlertCommand  string `json:"alert_command"`
}

type Config struct {
	Monitors []MonitorConf `json:"monitors"`
	Warnings []WarningConf `json:"warnings"`
}

func LoadConfig(configPath string) (Config, error) {
	config := Config{}
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	ext := path.Ext(configPath)
	if ext == ".json" {
		if err = json.Unmarshal(content, &config); err != nil {
			return config, err
		}
	} else if ext == ".yml" {
		if err = yaml.Unmarshal(content, &config); err != nil {
			return config, err
		}
	} else {
		return config, errors.New("not valid config")
	}
	return config, nil
}
