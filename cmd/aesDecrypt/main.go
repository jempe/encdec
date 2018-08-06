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
		fmt.Println("AES decrypter")
		fmt.Println("Decrypt the source folder and save it in the target folder")
		fmt.Println("Usage: aesdecrypt source_file target_folder")
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

			if strings.ToLower(filepath.Ext(source)) == ".enc" {
				data, err := ioutil.ReadFile(source)

				if err != nil {
					fmt.Println(err)
				} else {
					_, file := filepath.Split(source)

					targetFile := target + file

					if encdec.Exists(targetFile) {
						fmt.Println("file", targetFile, "already exists")
					} else {
						decrypted, err := encdec.Decrypt(data, []byte(key))

						err = ioutil.WriteFile(targetFile, decrypted, 0644)

						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println(source, "succesfully decrypted :", targetFile)
						}

					}
				}
			} else {
				fmt.Println("File", source, "is not encrypted")
			}
		}
	}
}
