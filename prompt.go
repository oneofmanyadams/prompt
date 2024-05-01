// Package prompt provides a simple and easy to use tool for creating
// command line prompts that end users can provide input to.
package prompt

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ask is the base function that performs the output of question q to the user
// and returns the answer provided by the user.
// If an error occurs, answer is an empty string.
func ask(q string, input io.Reader, output io.Writer) (string, error) {
	rdr := bufio.NewReader(input)

	output.Write([]byte(q + "\n"))
	output.Write([]byte("# "))

	raw_answer, read_error := rdr.ReadString('\n')

	if read_error != nil {
		return "", read_error
	}

	input_cleaner := strings.NewReplacer("\n", "")
	answer := strings.TrimSpace(input_cleaner.Replace(raw_answer))

	return answer, read_error
}

// Basic is designed to provide a single line way to get user input
// from the CLI.
func Basic(question string) (answer string, err error) {
	return ask(question, os.Stdin, os.Stdout)
}
