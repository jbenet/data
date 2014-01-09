package data

import (
	// "code.google.com/p/gcfg"
	"fmt"
	"github.com/gonuts/flag"
	"github.com/jbenet/commander"
	"os"
	"os/user"
	"strings"
)

// WARNING: the config format will be ini eventually. Go parsers
// don't currently allow writing (modifying) of files.
// Thus, for now, using yaml. Expect this to change.

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
	Run:  configCmd,
	Flag: *flag.NewFlagSet("data-config", flag.ExitOnError),
}

func init() {
	cmd_data_config.Flag.Bool("show", false, "show config file")
}

func configCmd(c *commander.Command, args []string) error {
	if c.Flag.Lookup("show").Value.Get().(bool) {
		return printConfig(&Config)
	}

	if len(args) == 0 {
		return fmt.Errorf("%s: requires <key> argument.", c.Name())
	}

	if len(args) == 1 {
		value, err := configGet(&Config, args[0])
		if err != nil {
			return err
		}

		pOut("%s\n", value)
		return nil
	}

	return configSet(&Config, args[0], args[1])
}

func printConfig(c *ConfigFormat) error {
	f, _ := NewConfigfile("")
	f.ConfigFormat = *c
	return f.Write(os.Stdout)
}

func configGet(c *ConfigFormat, key string) (string, error) {
	return "", NotImplementedError
}

func configSet(c *ConfigFormat, key string, value string) error {
	return NotImplementedError
}

var globalConfigFile = "~/.dataconfig"

type ConfigFormat struct {
	Index map[string]*struct {
		Url      string
		User     string
		Token    string
		Disabled bool ",omitempty"
	}
}

var Config ConfigFormat

// var DefaultConfigText = `[index "datadex.io:8080"]
// user =
// token =
// `
var DefaultConfigText = `index:
  datadex:
    url: http://datadex.io:8080
    user: ""
    token: ""
`

// Load config file on statup
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

func WriteConfigFile(filename string, fmt *ConfigFormat) error {
	// return gcfg.WriteFile(fmt, filename)

	f, _ := NewConfigfile(filename)
	f.ConfigFormat = *fmt
	return f.WriteFile()
}

func ReadConfigFile(filename string, fmt *ConfigFormat) error {
	// return gcfg.ReadFileInto(fmt, filename)

	f, err := NewConfigfile(filename)
	if err != nil {
		return err
	}

	*fmt = f.ConfigFormat
	return nil
}

// for use with YAML-based config
type Configfile struct {
	SerializedFile "-"
	ConfigFormat   ",inline"
}

func NewConfigfile(path string) (*Configfile, error) {
	f := &Configfile{SerializedFile: SerializedFile{Path: path}}
	f.SerializedFile.Format = f

	if len(path) > 0 {
		err := f.ReadFile()
		if err != nil {
			return f, err
		}
	}
	return f, nil
}
