// Package prompt is a simple tool for creating command line prompts.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Prompt struct {
	Question string
	Options map[string]string
	OptionOrder []string

	ExitKey string
	ExitValue string

	AnswerKey string
	AnswerValue string
	ValidAnswer bool

	InputFrom io.Reader
	OutputTo io.Writer

	errs []error

	QuickOptionCount int
}

func New(question string) (new_prompt Prompt) {
	new_prompt.Question = question
	new_prompt.Options = make(map[string]string)
	new_prompt.QuickOptionCount = 1
	new_prompt.InputFrom = os.Stdin
	new_prompt.OutputTo = os.Stdout
	return
}

func (p *Prompt) optionExists(option_key string, option_value string) (exists bool){
	for existing_key, existing_value := range p.Options {
		if option_key == existing_key {
			exists = true
		}
		if option_value == existing_value {
			exists = true
		}
	}
	return
}

func (p *Prompt) addOption(option_key string, option_value string) {
	p.Options[option_key] = option_value
	p.OptionOrder = append(p.OptionOrder, option_key)
}

func (p *Prompt) Option(option_key string, option_value string) (is_selected bool) {
	k := strings.TrimSpace(option_key)
	v := strings.TrimSpace(option_value)

	if !p.optionExists(k, v) {
		p.addOption(k, v)
	}

	is_selected = p.AnswerKey == k

	return
}

func (p *Prompt) O(option_value string) (is_selected bool) {
	option_key := strconv.Itoa(p.QuickOptionCount)
	k := strings.TrimSpace(option_key)
	v := strings.TrimSpace(option_value)

	if !p.optionExists(k, v) {
		p.QuickOptionCount = p.QuickOptionCount + 1
	}
	is_selected = p.Option(k, v)
	return
}

func QuickPrompt(question string, input_from io.Reader, output_to io.Writer) (answer string, errs []error) {
	rdr := bufio.NewReader(input_from)

	output_to.Write([]byte(question+"\n"))
	output_to.Write([]byte("#: ")) //+"\n" <---- Do I need this?

	raw_answer, read_error := rdr.ReadString('\n')

	if read_error != nil {
		errs = append(errs, read_error)
	}

	cleanup_input := strings.NewReplacer("\n", "")
	answer = strings.TrimSpace(cleanup_input.Replace(raw_answer))

	return
}

func (p *Prompt) optionsQuestion() (question string) {
	question = p.Question

	for _, key := range p.OptionOrder {
		question = fmt.Sprintf("%s \n %s %s", question, key, p.Options[key])
	}
	
	return
}

func(p* Prompt) answerInOptions(answer string) (exists bool, ans_key string) {
	for key, option := range p.Options {
		if option == answer {
			exists = true
			ans_key = key
			return
		}
		if key == answer {
			exists = true
			ans_key = key
			return
		}
	}
	return
}

func (p *Prompt) GetInput() {
	p.ValidAnswer = false
	question_string := p.optionsQuestion()
	answer, prompt_errors := QuickPrompt(question_string, p.InputFrom, p.OutputTo)

	if len(prompt_errors) > 0 {
		for _, new_err := range prompt_errors {
			p.errs = append(p.errs, new_err)
		} 
	}
	
	if len(p.Options) > 0 {
		if option_exists, existing_key := p.answerInOptions(answer); option_exists {
			p.AnswerKey = existing_key
			p.AnswerValue = p.Options[existing_key]
			p.ValidAnswer = true
		} else {
			p.AnswerKey = ""
			p.AnswerValue = answer
			p.ValidAnswer = false
		}
	} else {
		p.AnswerKey = ""
		p.AnswerValue = answer
		p.ValidAnswer = false
	}
}