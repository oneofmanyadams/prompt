package prompt

import (
	"os"
)

// Basic is designed to provide a single line way to get user input
// from the CLI.
func Basic(question string) (answer string, err error) {
	return ask(question, os.Stdin, os.Stdout)
}
