package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Lang string `json:"lang"`
}

const configPath = ".leetcode.json"

func LoadConfig() Config {
	var c Config
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return c
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		return c
	}
	return c
}

func GetLang() string {
	c := LoadConfig()
	return c.Lang
}
