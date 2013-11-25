package main

const VERSION = "0.0.2"

func VersionCmd([]string) error {
	DOut("data version %s\n", VERSION)
	return nil
}
