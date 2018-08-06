package main

import (
	"bitbucket.org/suciokastro/encdec"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	} else if len(os.Args) < 3 {
		fmt.Println("AES encrypter")
		fmt.Println("Encrypt the source folder and save it in the target folder")
		fmt.Println("Usage: aesencrypt source_file target_folder")
	} else {
		source := os.Args[1]
		target := os.Args[2]

		if !encdec.Exists(source) {
			fmt.Println("source file", source, "doesn't exist")
		} else if !encdec.Exists(target) {
			fmt.Println("target folder", target, "doesn't exist")
		} else if encdec.IsDir(source) {
			fmt.Println("source", source, "can't be a folder")
		} else if !encdec.IsDir(target) {
			fmt.Println("target", target, "is not a folder")
		} else {
			target = strings.TrimSuffix(target, "/") + "/"

			if strings.ToLower(filepath.Ext(source)) != ".enc" {
				data, err := ioutil.ReadFile(source)

				if err != nil {
					fmt.Println(err)
				} else {
					md5Sum := encdec.Md5Sum(data)

					encryptedName := md5Sum + ".enc"

					targetFile := target + encryptedName

					if encdec.Exists(targetFile) {
						fmt.Println("file", source, "is already encrypted :", targetFile)
					} else {
						encrypted, err := encdec.Encrypt(data, []byte(key))

						err = ioutil.WriteFile(targetFile, encrypted, 0644)

						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println(source, "succesfully encrypted :", targetFile)
						}

					}
				}
			} else {
				fmt.Println("Can't encrypt encrypted file", source)
			}
		}
	}
}
