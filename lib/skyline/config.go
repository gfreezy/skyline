package skyline

import (
	"encoding/json"
	"io/ioutil"
)

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

func LoadConfig(path string) (Config, error) {
	config := Config{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return config, err
	}
	return config, nil
}
