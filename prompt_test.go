package prompt

import (
	"fmt"
	"testing"
	"strings"
	"time"
	"os"
	"io/ioutil"
)

//////////////////////////////////////////////////////////////////
// Examples
//////////////////////////////////////////////////////////////////

func Example() {
	prmpt := NewPrompt("What is your favorite color?")
	prmpt.AddOption("1", "Red")
	prmpt.AddOption("2", "Blue")
	prmpt.AddOption("3", "Green")
	answer, blunders := prmpt.PromptRequireOption()
	fmt.Print(answer)
	fmt.Print(blunders)
}

func ExampleQuickPrompt() {
	answer, errs := QuickPrompt("Who are you?", os.Stdin, os.Stdout)
	fmt.Print(answer)
	fmt.Print(errs)
}

func ExampleNewPrompt() {
	prmpt := NewPrompt("What is your favorite color?")
	answer, blunders := prmpt.PromptUser()
	fmt.Print(answer)
	fmt.Print(blunders)
}

func ExamplePrompt_AddOption() {
	prmpt := NewPrompt("What is your favorite color?")
	prmpt.AddOption("1", "Red")
	prmpt.AddOption("2", "Blue")
	prmpt.AddOption("3", "Green")
}


func ExamplePrompt_PromptUser() {
	prmpt := NewPrompt("What is your favorite color?")
	prmpt.AddOption("1", "Red")
	prmpt.AddOption("2", "Blue")
	prmpt.AddOption("3", "Green")
	answer, blunders := prmpt.PromptUser()
	fmt.Print(answer)
	fmt.Print(blunders)
}

func ExamplePrompt_PromptRequireOption() {
	prmpt := NewPrompt("What is your favorite color?")
	prmpt.AddOption("1", "Red")
	prmpt.AddOption("2", "Blue")
	prmpt.AddOption("3", "Green")
	answer, blunders := prmpt.PromptRequireOption()
	fmt.Print(answer)
	fmt.Print(blunders)
}

func ExamplePrompt_GetInputFrom() {
	prmpt := NewPrompt("What is your favorite color?")
	prmpt.GetInputFrom(strings.NewReader("Forced Answer\n"))
}

func ExamplePrompt_SendOutputTo() {
	prmpt := NewPrompt("What is your favorite color?")
	prmpt.SendOutputTo(ioutil.Discard)
}

//////////////////////////////////////////////////////////////////
// Initialization Tests
//////////////////////////////////////////////////////////////////

func TestNewPrompt(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")
	if prmpt.Question != "TestQuestion" {
		t.Errorf("Prompt instance Question field did not match provided question.")
	}
	if len(prmpt.Options) != 0 {
		t.Errorf("Prompt instance initialized with pre-existing options.")
	}
	if len(prmpt.Blunders.Blunders) > 0 {
		t.Errorf("Prompt instance initializing with blunders.")
	}
}

//////////////////////////////////////////////////////////////////
// Option Tests
//////////////////////////////////////////////////////////////////

func TestAddOption_single_add(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("1", "one")
	if !added || len(prmpt.Blunders.Blunders) > 0 {
		blunder_string := prmpt.Blunders.BlunderSliceAsString(prmpt.Blunders.Blunders)
		t.Errorf("Unable to add single option. Failed with blunders: \n"+blunder_string)
	}
}

func TestAddOption_multiple_add(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added1 := prmpt.AddOption("1", "one")
	added2 := prmpt.AddOption("2", "two")
	added3 := prmpt.AddOption("3", "three")
	added4 := prmpt.AddOption("4", "four")
	if !added1 || !added2 || !added3 || !added4 || len(prmpt.Blunders.Blunders) > 0 {
		blunder_string := prmpt.Blunders.BlunderSliceAsString(prmpt.Blunders.Blunders)
		t.Errorf("Unable to add multiple options. Failed with blunders: \n"+blunder_string)
	}
}

func TestAddOption_add_existing_key(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("1", "one")
	added1 := prmpt.AddOption("1", "not_one")
	// makes sure added1 is false and that a blunder was recorded.
	// should split into different blocks for added1 and Blunders.Reported len()
	if !added || added1 || len(prmpt.Blunders.Blunders) < 1 {
		t.Errorf("Was able to add 2 options with same key. (should not be possible)")
	}
}

func TestAddOption_add_existing_value(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("1", "one")
	added1 := prmpt.AddOption("not_1", "one")
	// makes sure added1 is false and that a blunder was recorded.
	// should split into different blocks for added1 and Blunders.Reported len()
	if !added || added1 || len(prmpt.Blunders.Blunders) < 1 {
		t.Errorf("Was able to add options with same value. (should not be possible)")
	}
}

func TestAddOption_add_empty_key(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("", "one")
	if added || len(prmpt.Blunders.Blunders) < 1 {
		t.Errorf("Was able to add an option with no key. (should not be possible)")
	}
}

func TestAddOption_add_empty_value(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("1", "")
	if added || len(prmpt.Blunders.Blunders) < 1 {
		t.Errorf("Was able to add an option with no value. (should not be possible)")
	}
}

//////////////////////////////////////////////////////////////////
// User Prompting Tests
//////////////////////////////////////////////////////////////////

func TestQuickPrompt(t *testing.T) {
	answer := "Answer"
	result, err := QuickPrompt("Question", strings.NewReader(answer+"\n"), ioutil.Discard)

	if result != answer {
		t.Errorf("Did not return the information the way the user provided it.")
	}
	if err != nil {
		t.Errorf("Failed with Error: "+ err.Error())
	}
}

func TestPromptUser_without_options(t *testing.T) {
	answer := "Answer"
	prmpt := NewPrompt("TestQuestion")
	prmpt.GetInputFrom(strings.NewReader(answer+"\n"))
	prmpt.SendOutputTo(ioutil.Discard)
	result, blndr := prmpt.PromptUser()

	if result != answer {
		t.Errorf("Result did not match user input. Got: %s, Expected: %s", result, answer)
	}
	if prmpt.Answer != answer {
		t.Errorf("Answer was not saved in Prompt.Answer . Got: %s, Expected: %s", prmpt.Answer, answer)
	}
	if blndr.Message != "" {
		t.Errorf("Got Blunder: "+blndr.Error())
	}

}

func TestPromptUser_using_keys(t *testing.T) {

	type test_struct struct {
		k string
		v string
	}
	tables := []test_struct{
			{"1", "one"},
			{"2", "two"},
			{"3", "three"}}


	prmpt := NewPrompt("TestQuestion")

	for _, tester := range tables {
		prmpt.AddOption(tester.k, tester.v)
	}

	for _, answer := range tables {
		prmpt.GetInputFrom(strings.NewReader(answer.k+"\n"))
		prmpt.SendOutputTo(ioutil.Discard)
		result, blndr := prmpt.PromptUser()
		if result != answer.v {
			t.Errorf("Did not get expected result based on answer. got: %s, expected: %s.", result, answer.v)
		}
		if blndr.Message != "" {
			t.Errorf("Got Blunder: "+blndr.Error())
		}
		if prmpt.Answer != answer.v {
			t.Errorf("Answer was not saved in Prompt.Answer . Got: %s, Expected: %s", prmpt.Answer, answer.v)
		}	
		if len(prmpt.Blunders.Blunders) > 0 {
			blunder_string := prmpt.Blunders.BlunderSliceAsString(prmpt.Blunders.Blunders)
			t.Errorf("Found Blunders: \n"+blunder_string)
		}
	}

}

func TestPromptUser_using_values(t *testing.T) {

	type test_struct struct {
		k string
		v string
	}
	tables := []test_struct{
			{"1", "one"},
			{"2", "two"},
			{"3", "three"}}


	prmpt := NewPrompt("TestQuestion")

	for _, tester := range tables {
		prmpt.AddOption(tester.k, tester.v)
	}

	for _, answer := range tables {
		prmpt.GetInputFrom(strings.NewReader(answer.v+"\n"))
		prmpt.SendOutputTo(ioutil.Discard)
		result, blndr := prmpt.PromptUser()
		if result != answer.v {
			t.Errorf("Did not get expected result based on answer. got: %s, expected: %s.", result, answer.v)
		}
		if prmpt.Answer != answer.v {
			t.Errorf("Answer was not saved in Prompt.Answer . Got: %s, Expected: %s", prmpt.Answer, answer.v)
		}	
		if blndr.Message != "" {
			t.Errorf("Got Blunder: "+blndr.Error())
		}
		if len(prmpt.Blunders.Blunders) > 0 {
			blunder_string := prmpt.Blunders.BlunderSliceAsString(prmpt.Blunders.Blunders)
			t.Errorf("Found Blunders: \n"+blunder_string)
		}
	}

}

func TestPromptUser_wrong_option(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")
	prmpt.AddOption("1", "one")
	prmpt.GetInputFrom(strings.NewReader("4\n"))
	prmpt.SendOutputTo(ioutil.Discard)

	result, blndr := prmpt.PromptUser()

	if result != "" {
		t.Errorf("Failed to set result to empty string. Returned \"%s\" instead", result)
	}
	if prmpt.Answer != "" {
		t.Errorf("Malformed answer was saved in Prompt.Answer . Got: %s, Expected: %s", prmpt.Answer, "")
	}
	if blndr.Message == "" {
		t.Errorf("Failed to record a blunder")
	}

}

// This whole func seems a bit wonky
func TestPromptRequireOption(t *testing.T) {
	p := NewPrompt("Test Question")
	prmpt := &p
	prmpt.AddOption("1", "one")
	prmpt.GetInputFrom(strings.NewReader("4\n"))
	prmpt.SendOutputTo(ioutil.Discard)

	// This is kinda weird. Should review this.
	go func(p *Prompt) {
		time.Sleep(1 * time.Millisecond)
		prmpt.GetInputFrom(strings.NewReader("1\n"))
	}(prmpt)

	result, _ := prmpt.PromptRequireOption()

	if result != "one" {
		t.Errorf("Never recovered")
	}
	if prmpt.Answer != "one" {
		t.Errorf("Answer was not saved in Prompt.Answer . Got: %s, Expected: %s", prmpt.Answer, "one")
	}
	if len(prmpt.Blunders.Blunders) < 1 {
		t.Errorf("Failed to record a blunder")
	}
}

func TestOptionsQuestion(t *testing.T) {
	options := [][2]string{{"1", "one"},{"2", "two"}}	
	
	prmpt := NewPrompt("Test Question")
	prmpt.SendOutputTo(ioutil.Discard)

	option_string_manual := prmpt.Question

	for _, option := range options {
		prmpt.AddOption(option[0], option[1])
		option_string_manual = option_string_manual + " \n "+option[0]+" "+option[1]
	}

	options_string_generated := prmpt.optionsQuestion()

	if options_string_generated != option_string_manual {
		t.Errorf("Generated option string does not match manually built one.")
	}
}

func TestAnswerInOptions(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")
	prmpt.AddOption("1", "one")
	prmpt.AddOption("2", "two")
	prmpt.AddOption("3", "three")

	if !prmpt.answerInOptions("1") {
		t.Errorf("Existing key answer did not return true.")
	}
	if !prmpt.answerInOptions("two") {
		t.Errorf("Existing value answer did not return true.")
	}
	if prmpt.answerInOptions("four") {
		t.Errorf("Non-Existing answer returned true.")
	}

}

func TestGetInputFrom(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	if prmpt.InputFrom != os.Stdin {
		t.Errorf("InputFrom not defaulting to os.Stdin.")
	}

	prmpt.GetInputFrom(strings.NewReader("1\n"))
	prmpt.SendOutputTo(ioutil.Discard)


	if prmpt.InputFrom == os.Stdin {
		t.Errorf("GetInputFrom is not setting to provided value.")
	}
}