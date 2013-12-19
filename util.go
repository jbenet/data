package data

import (
	"fmt"
	"os"
	"os/exec"
	"unicode"
)

var Debug bool

// Shorthand printing functions.
func pErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func pOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func dErr(format string, a ...interface{}) {
	if Debug {
		pErr(format, a...)
	}
}

func dOut(format string, a ...interface{}) {
	if Debug {
		pOut(format, a...)
	}
}

// Checks whether string is a hash (sha1)
func isHash(hash string) bool {
	if len(hash) != 40 {
		return false
	}

	for _, r := range hash {
		if !unicode.Is(unicode.ASCII_Hex_Digit, r) {
			return false
		}
	}

	return true
}

func shortHash(hash string) string {
	return hash[:7]
}

func copyFile(src string, dst string) error {
	cmd := exec.Command("cp", src, dst)
	return cmd.Run()
}
