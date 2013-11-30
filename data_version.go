package data

const Version = "0.0.2"

func versionCmd([]string) error {
	pOut("data version %s\n", Version)
	return nil
}
