package prompt

import (
	"io"
	"os"
)

// OptionPrompt is used for creating a set of options to present to a user.
// This struct should not be used directly, but instead created by calling
// the New() function.
type OptionPrompt struct {
	Question string
	Options  []string
	Answer   string
	Ask      bool

	Input  io.Reader
	Output io.Reader

	UserInput string
}

////////////////////////////////////////////////////////////////////////////////
// Constructor
////////////////////////////////////////////////////////////////////////////////

// New should be the only method for creating an OptionPrompt instance.
func New(question string) OptionPrompt {
	return OptionPrompt{
		Question: question,
		Ask:      true,
		Input:    os.Stdin,
		Output:   os.Stdout}
}

////////////////////////////////////////////////////////////////////////////////
// Public Functions
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Private Functions
////////////////////////////////////////////////////////////////////////////////

func (s *OptionPrompt) optionExists(option_name string) (int, bool) {
	for k, v := range s.Options {
		if option_name == v {
			return k, true
		}
	}
	return len(s.Options), false
}

func (s *OptionPrompt) addOption(option_name string) int {
	s.Options = append(s.Options, option_name)
	return len(s.Options) - 1
}
