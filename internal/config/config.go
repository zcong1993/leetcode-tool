package config

import (
	"encoding/json"
	"io/ioutil"
)

const defaultCookie = "alicfw=1089899001%7C2081167297%7C1328233593%7C1328234352; alicfw_gfver=v1.200309.1"

type Config struct {
	Lang   string `json:"lang"`
	Cookie string `json:"cookie"`
}

const configPath = ".leetcode.json"

func NewConfig() *Config {
	c := loadConfig()
	if c.Cookie == "" {
		c.Cookie = defaultCookie
	}
	return &c
}

func loadConfig() Config {
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
