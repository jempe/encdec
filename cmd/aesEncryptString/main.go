package main

import (
	"bitbucket.org/suciokastro/encdec"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

var configPath string
var key string

func main() {
	if runtime.GOOS == "windows" {
		configPath = os.Getenv("APPDATA") + "/aesencdec"
	} else {
		configPath = os.Getenv("HOME") + "/.aesencdec"
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	viper.ReadInConfig()

	key = viper.GetString("key")

	if key == "" {
		fmt.Println("Please set a key in the configuration folder", configPath)
	} else if len(os.Args) < 2 {
		fmt.Println("AES string encrypter")
		fmt.Println("Encrypt a string")
		fmt.Println("Usage: aesEncryptString string")
	} else {
		out, err := encdec.EncryptString(os.Args[1], []byte(key))

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(out)
		}
	}
}
