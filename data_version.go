package data

const VERSION = "0.0.2"

func VersionCmd([]string) error {
	Out("data version %s\n", VERSION)
	return nil
}
