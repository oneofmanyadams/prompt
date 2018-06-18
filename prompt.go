// Package prompt is a simple tool for creating command line prompts.
package prompt

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"blunders"
)

// Prompt is the main type for the promp package.
//  - Question is the question that the user will be prompted with.
//  - Options is a map of an option_id and it's associated name.
//  - Order keeps track of what order options are added in and makes sure they are always displayed in that order.
//  - Answer provides a place to store the most recent answer provided by the user.
//  - InputFrom determines where to get user input from. (os.Stdin is the usual setting)
//  - Blunders the implementation of a custom package that expands error recording and handling.
type Prompt struct {
	Question string
	Options map[string]string
	Order []string
	Answer string
	InputFrom io.Reader
	Blunders blunders.Blunders
}

// NewPrompt creates a new instace of Prompt.
// Sets the Question and initializes the Options map.
// Sets InputFrom to os.Stdin
// Initializes the Blunders instance.
// This is where all Blunder Codes are created. 
func NewPrompt(question string) (p Prompt) {
	p.Question = question
	p.Options = make(map[string]string)
	p.Blunders = blunders.NewBlunders("Prompt")
	p.Blunders.AddCode(1, "OptionAddError")
	p.Blunders.AddCode(2, "UserInputError")
	p.InputFrom = os.Stdin
	return
}

//////////////////////////////////////////////////////////////////
// Option Functions
//////////////////////////////////////////////////////////////////

// AddOption is the correct way to add an option to a prompt.
// Adds key and question to Options map[key]question.
// Adds the key to Order []key.
// If the key or question provided already exist, a FATAL blunder is reported.
// An empty string for key or question also results in a FATAL blunder being reported. 
func (p *Prompt) AddOption(key string, question string) (added bool) {
	added = true
	
	if key == "" {
		p.Blunders.NewFatal(1, "Empty string provided for key.")
		added = false
	}

	if question == "" {
		p.Blunders.NewFatal(1, "Empty string provided for question.")
		added = false
	}
	
	for existing_key, existing_question := range p.Options {
		if existing_key == key {
			p.Blunders.NewFatal(1, fmt.Sprintf("Attempted to add already existing option key \"%s\".", key))
			added = false
		}
		if existing_question == question {
			p.Blunders.NewFatal(1, fmt.Sprintf("Attempted to add already existing option question \"%s\".", question))
			added = false
		}
	}

	if added {
		p.Options[key] = question
		p.Order = append(p.Order, key)
	}
	
	return
}

//////////////////////////////////////////////////////////////////
// User Prompting Functions
//////////////////////////////////////////////////////////////////

// QuickPrompt is the most basic function that will promp a user.
// It is the only function that can prompt a user WITHOUT being called
// as a method of a Prompt instance.
// 2nd argument is an io.Reader to where the input is coming from (typically os.Stdin).
// It uses an error type as it's 2nd return value instead of a blunder for simplicities sake.
// The user input is stopped being captured at the first detection of a newline \n.
func QuickPrompt(question string, input_from io.Reader) (answer string, err error) {
	rdr := bufio.NewReader(input_from)

	fmt.Println(question)
	fmt.Print("#:")

	raw_answer, read_error := rdr.ReadString('\n')

	if read_error != nil {
		err = read_error
	}

	cleanup_input := strings.NewReplacer("\n", "")
	answer = strings.TrimSpace(cleanup_input.Replace(raw_answer))

	return
}

// PromptUser is default way to prompt a user.
// If the prompt instance has any options, they are automatically loaded
// and presented to the user.
// If options are presented, the function will return a non-fatal blunder
// if the answer provided is NOT an option.
// The answer is considered valid if it matches either an option key or an option question.
// If the provided answer matches a key, answer will automatically be converted to the question value.
func (p *Prompt) PromptUser() (answer string, blndr blunders.Blunder) {

	question_string := p.optionsQuestion()
	var prompt_error error
	answer, prompt_error = QuickPrompt(question_string, p.InputFrom)

	if prompt_error != nil {
		blndr = p.Blunders.New(2, prompt_error.Error())
		answer = ""
	}
	
	if len(p.Options) > 0 && blndr.Fatal == false {
		if p.answerInOptions(answer) == false {
			blndr = p.Blunders.New(2, fmt.Sprintf("Option provided (%s) does not exist.", answer))
			answer = ""
		} else {
			for key := range p.Options {
				if key == answer {
					answer = p.Options[key]
				}
			}
		}
	}

	return
}

// PromptRequireOption is very similar to PromptUser but will continue to
// prompt the user until a valid option is entered.
func (p *Prompt) PromptRequireOption() (answer string, blndr blunders.Blunder) {
	answer, blndr = p.PromptUser()
	for ; blndr.Code != 0 ; {
		fmt.Println("!!"+blndr.Message)
		answer, blndr = p.PromptUser()
	}
	return
}

//////////////////////////////////////////////////////////////////
// Helper Functions
//////////////////////////////////////////////////////////////////

// optionsQuestion combines all options into 1 string that will display
// one option per line in the command line.
// The combined option string is added to the back of the qustion and
// returned as one string.
func (p *Prompt) optionsQuestion() (question string) {
	question = p.Question

	for _, key := range p.Order {
		question = fmt.Sprintf("%s \n %s %s", question, key, p.Options[key])
	}
	
	return
}

// answerInOptions takes an answer string and checks to see if an identical
// string exists as either a option key or an option question.
// Returns true if a match is found and false if no match is found.
func(p* Prompt) answerInOptions(answer string) (exists bool) {
	for key, option := range p.Options {
		if option == answer {
			exists = true
			return
		}
		if key == answer {
			exists = true
			return
		}
	}
	return
}

//////////////////////////////////////////////////////////////////
// Utility Functions
//////////////////////////////////////////////////////////////////

// GetInputFrom provides and easy way to change where a Prompt instance
// will get it's input from. (defaults to os.Stdin on initialization in NewPrompt(string))
func (p* Prompt) GetInputFrom(input io.Reader) {
	p.InputFrom = input
}