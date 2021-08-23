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

func GetCookie() string {
	c := LoadConfig()
	ck := c.Cookie
	if ck == "" {
		return defaultCookie
	}
	return ck
}
