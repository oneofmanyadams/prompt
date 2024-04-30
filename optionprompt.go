package prompt

import (
	"io"
)

type OptionPrompt struct {
	Question string
	Options  []string
	Answer   string
	Ask      bool

	Input  io.Reader
	Output io.Reader

	UserInput string
}
