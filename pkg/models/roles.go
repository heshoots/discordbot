package models

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Role struct {
	Id     int64  `json:"ID,omitempty", yaml: "ID"`
	Name   string `json:"name", yaml: "name"`
	RoleID string `json:"roleID", yaml: "roleID"`
}

type Roles struct {
	Roles []*Role `json:"roles"`
}

func save() {
	viper.SafeWriteConfig()
}

func YamlRoles() ([]*Role, error) {
	var roles Roles
	err := viper.Unmarshal(&roles)
	return roles.Roles, err
}

func YamlRole(name string) (*Role, error) {
	var roles Roles
	err := viper.Unmarshal(&roles)
	if err != nil {
		return nil, err
	}
	for _, i := range roles.Roles {
		if strings.ToUpper(i.Name) == strings.ToUpper(name) {
			return i, nil
		}
	}
	return nil, errors.New("Couldn't find role")
}
