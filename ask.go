package prompt

import (
	"bufio"
	"io"
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
