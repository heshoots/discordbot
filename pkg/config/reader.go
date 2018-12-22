package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/heshoots/discordbot/pkg/models"
)

var config models.Roles
var configLoc string

func ReadConfigFile(path string) models.Roles {
	configLoc = path
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return ReadConfig(string(data))
}

func AddRole(role *models.Role) error {
	config.Roles = append(config.Roles, role)
	return WriteConfigFile()
}

func ReadConfig(configStr string) models.Roles {
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func GetConfig() models.Roles {
	return config
}

func WriteConfigFile() error {
	err := ioutil.WriteFile(configLoc, []byte(WriteConfig()), 0644)
	if err != nil {
		return err
	}
	return nil
}

func WriteConfig() string {
	configstr, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return ""
	}
	return string(configstr)
}
