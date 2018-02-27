package models

import (
	"log"
)

type Role struct {
	Id     int64
	Name   string
	RoleID string
}

func CreateRole(role *Role) error {
	err := db.Insert(role)
	if err != nil {
		return err
	}
	return nil
}

func GetRole(name string) (*Role, error) {
	role := Role{}
	err := db.Model(&role).Where("name ILIKE ?", name).Select()
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func DeleteRole(name string) error {
	_, err := db.Model(&Role{}).Where("name ILIKE ?", name).Delete()
	return err
}

func AddRoleCall(roleName string) {
	role := Role{}
	err := db.Model(&role).Where("name = ?", roleName).Select()
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println(role)
	_, err = db.Exec("INSERT into rolecall (timestamp, role) VALUES (NOW(), ?)", role.Id)
	if err != nil {
		log.Panic(err)
		return
	}
}

func GetRoles() ([]Role, error) {
	var roles []Role
	err := db.Model(&roles).Select()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
