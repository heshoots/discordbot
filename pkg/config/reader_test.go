package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	config := ReadConfig(`{"roles":[{"roleID":"12","name":"testrole"}]}`)
	if config.Roles[0].RoleID != "12" {
		t.Errorf("failed to get roles" + config.Roles[0].RoleID + " != 12")
	}
	if config.Roles[0].Name != "testrole" {
		t.Errorf("failed to get roles" + config.Roles[0].Name + " != testrole")
	}
}

func TestWriteConfig(t *testing.T) {
	teststr := `{
  "roles":[{"name":"testrole","roleID":"12"}]}`
	ReadConfig(teststr)
	outstr := WriteConfig()
	if outstr != teststr {
		t.Errorf("Did not return same configuration passed in: " + outstr)
	}
}
