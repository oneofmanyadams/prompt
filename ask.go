package prompt

import (
	"bufio"
	"io"
	"strings"
)

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
