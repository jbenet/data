package data

import (
	"code.google.com/p/gcfg"
	"os"
	"os/user"
	"strings"
)

var globalConfigFile = "~/.dataconfig"

type ConfigFormat struct {
	Index map[string]*struct {
		User     string
		Token    string
		Disabled bool
	}
}

var Config ConfigFormat

var DefaultConfigText = `[index "datadex.io:8080"]
user =
token =
`

func init() {

	// expand ~/
	usr, err := user.Current()
	if err != nil {
		panic("error: user context.")
	}
	dir := usr.HomeDir + "/"
	globalConfigFile = strings.Replace(globalConfigFile, "~/", dir, 1)

	// install config if doesn't exist
	if _, err := os.Stat(globalConfigFile); os.IsNotExist(err) {
		err := WriteConfigFileText(globalConfigFile, DefaultConfigText)
		if err != nil {
			panic("error: failed to write config " + globalConfigFile +
				". " + err.Error())
		}
	}

	// load config
	err = ReadConfigFile(globalConfigFile, &Config)
	if err != nil {
		panic("error: failed to load config " + globalConfigFile +
			". " + err.Error())
	}
}

func WriteConfigFileText(filename string, text string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(text))
	return err
}

func ReadConfigFile(filename string, fmt *ConfigFormat) error {
	return gcfg.ReadFileInto(fmt, filename)
}
