// Package prompt is a simple tool for creating command line prompts.
package prompt

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"blunders"
)

// Prompt is the main type for the promp package.
//  - Question is the question that the user will be prompted with.
//  - Options is a map of an option_id and it's associated name.
//  - Order keeps track of what order options are added in and makes sure they are always displayed in that order.
//  - Answer provides a place to store the most recent answer provided by the user.
//  - Blunders the implementation of a custom package that expands error recording and handling.
type Prompt struct {
	Question string
	Options map[string]string
	Order []string
	Answer string
	Blunders blunders.Blunders
}

// NewPrompt creates a new instace of Prompt.
// Sets the Question and initializes the Options map.
// Initializes the Blunders instance.
// This is where all Blunder Codes are created. 
func NewPrompt(question string) (p Prompt) {
	p.Question = question
	p.Options = make(map[string]string)
	p.Blunders = blunders.NewBlunders("Prompt")
	p.Blunders.AddCode(1, "OptionAddError")
	p.Blunders.AddCode(2, "UserInputError")
	return
}

// AddOption is the correct way to add an option to a prompt.
// Adds key and question to Options map[key]question.
// Adds the key to Order []key.
// If the key or question provided already exist, a FATAL blunder is reported.
// An empty string key or question also results in a FATAL blunder being reported. 
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

func PromptUser(question string) (answer string, err error) {
	rdr := bufio.NewReader(os.Stdin)

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

func (p *Prompt) Prompt() (answer string, blndr blunders.Blunder) {

	question_string := p.optionsQuestion()
	var prompt_error error
	answer, prompt_error = PromptUser(question_string)

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

func (p *Prompt) PromptForOptions() (answer string, blndr blunders.Blunder) {
	answer, blndr = p.Prompt()
	for ; blndr.Code != 0 ; {
		fmt.Println("!!"+blndr.Message)
		answer, blndr = p.Prompt()
	}
	return
}

func (p *Prompt) optionsQuestion() (question string) {
	question = p.Question

	for _, key := range p.Order {
		question = fmt.Sprintf("%s \n %s %s", question, key, p.Options[key])
	}
	
	return
}

func(p* Prompt) answerInOptions(answer string) (exists bool) {
	for key, option := range p.Options {
		if option == answer {
			exists = true
			return
		}
		if key == answer {
			exists = true
			answer = option
			return
		}
	}
	return
}
