package config

import (
	"bytes"
	"embed"
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

//go:embed *.yaml
var configs embed.FS

func InitConfig() {
	env := os.Getenv("ENV")
	vp := viper.New()
	configStream, err := configs.ReadFile("application." + env + ".yaml")
	if err != nil {
		panic(err)
	}
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(bytes.NewReader(configStream))
	if err != nil {
		panic(err)
	}
	err = vp.UnmarshalKey("app", &App)
	if err != nil {
		panic(err)
	}
	App.AESKEY, err = hex.DecodeString(App.AESKeyStr)
	if err != nil {
		panic(err)
	}
	err = vp.UnmarshalKey("database", &DB)
	if err != nil {
		panic(err)
	}
	password, err := os.ReadFile(DB.Master.PasswordFilePath)
	if err != nil {
		panic(fmt.Errorf("filepath is %s. err is %s", DB.Master.PasswordFilePath, err))
	}
	DB.Master.Dsn = fmt.Sprintf(DB.Master.Dsn, strings.TrimSpace(string(password)))

	err = vp.UnmarshalKey("mq", &MQ)
	if err != nil {
		panic(err)
	}
	mqPassword, err := os.ReadFile(MQ.PasswordFilePath)
	if err != nil {
		panic(fmt.Errorf("filepath is %s. err is %s", MQ.PasswordFilePath, err))
	}
	MQ.Url = fmt.Sprintf(MQ.Url, strings.TrimSpace(string(mqPassword)))
}
