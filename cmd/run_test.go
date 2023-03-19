package cmd

import (
	"fmt"
	"github.com/klzwii/mirai-go/util"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	util.ConfigFile = fmt.Sprintf("%s/%s", path, "mirai_config_test.yaml")
	data, err := loadConfig()
	assert.Nil(t, err)
	botConfig := data["bot"].(map[string]interface{})
	assert.Equal(t, 1234567, botConfig["qq"])
	assert.Equal(t, 23456, botConfig["verify_key"])
	assert.Equal(t, 6789, botConfig["port"])
	fmt.Println(data)
}
