package prompt

import (
	"io"
	"os"
	"strconv"
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

func (s *OptionPrompt) Option(option_name string) bool {
	option_key, option_exists := s.optionExists(option_name)
	if option_exists == false {
		option_key = s.addOption(option_name)
	}

	// Convert key to string to compare to user input.
	if s.UserInput == strconv.Itoa(option_key) {
		s.Answer = option_name
		s.Ask = false
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// Private Functions
////////////////////////////////////////////////////////////////////////////////

// optionExists determines if the provided option_name already exists in s.Options.
// If it does the first returned calue is the key and the 2nd return value is true.
// If it doesn't the next available key is returned and the 2nd return value is false.
func (s *OptionPrompt) optionExists(option_name string) (int, bool) {
	for k, v := range s.Options {
		if option_name == v {
			return k, true
		}
	}
	return len(s.Options), false
}

// addOption adds an option to s.Options and returns it's key.
func (s *OptionPrompt) addOption(option_name string) int {
	s.Options = append(s.Options, option_name)
	return len(s.Options) - 1
}
