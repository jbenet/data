package data

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/jbenet/commander"
	"os"
	"os/user"
	"strings"
)

var cmd_data_config = &commander.Command{
	UsageLine: "config <command> <key> [<value>]",
	Short:     "Manage data configuration.",
	Long: `data config - Manage data configuration.

    Usage:

      data config <key> [<value>]

    Get or set configuration option values.
    If <value> argument is not provided, print <key> value, and exit.
    If <value> argument is provided, set <key> to <value>, and exit.

      # sets foo.bar = buzz
      > data config foo.bar baz

      # gets foo.bar
      > data config foo.bar
      baz

    Config options are stored in the user's configuration file (~/.dataconfig).
    This file is formatted like .gitconfig (INI style), and uses the gcfg parser.
  `,
	Run: configCmd,
}

func configCmd(c *commander.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%s: requires <key> argument.", c.Name())
	}

	if len(args) == 1 {
		value, err := configGet(args[0])
		if err != nil {
			return err
		}

		pOut("%s\n", value)
		return nil
	}

	return configSet(args[0], args[1])
}

func configGet(key string) (string, error) {
	return "", NotImplementedError
}

func configSet(key string, value string) error {
	return NotImplementedError
}

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
