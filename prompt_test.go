package prompt

import "testing"


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
	if len(prmpt.Blunders.Reported) > 0 {
		t.Errorf("Prompt instance initializing with blunders.")
	}
}

//////////////////////////////////////////////////////////////////
// Option Tests
//////////////////////////////////////////////////////////////////

func TestAddOption_single_add(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("1", "one")
	if !added || len(prmpt.Blunders.Reported) > 0 {
		blunder_string := prmpt.Blunders.BlunderListToLogString(prmpt.Blunders.Reported)
		t.Errorf("Unable to add single option. Failed with blunders: \n"+blunder_string)
	}
}

func TestAddOption_multiple_add(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added1 := prmpt.AddOption("1", "one")
	added2 := prmpt.AddOption("2", "two")
	added3 := prmpt.AddOption("3", "three")
	added4 := prmpt.AddOption("4", "four")
	if !added1 || !added2 || !added3 || !added4 || len(prmpt.Blunders.Reported) > 0 {
		blunder_string := prmpt.Blunders.BlunderListToLogString(prmpt.Blunders.Reported)
		t.Errorf("Unable to add multiple options. Failed with blunders: \n"+blunder_string)
	}
}

func TestAddOption_add_existing_key(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("1", "one")
	added1 := prmpt.AddOption("1", "not_one")
	// makes sure added1 is false and that a blunder was recorded.
	// should split into different blocks for added1 and Blunders.Reported len()
	if !added || added1 || len(prmpt.Blunders.Reported) < 1 {
		t.Errorf("Was able to add 2 options with same key. (should not be possible)")
	}
}

func TestAddOption_add_existing_value(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("1", "one")
	added1 := prmpt.AddOption("not_1", "one")
	// makes sure added1 is false and that a blunder was recorded.
	// should split into different blocks for added1 and Blunders.Reported len()
	if !added || added1 || len(prmpt.Blunders.Reported) < 1 {
		t.Errorf("Was able to add options with same value. (should not be possible)")
	}
}

func TestAddOption_add_empty_key(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("", "one")
	if added || len(prmpt.Blunders.Reported) < 1 {
		t.Errorf("Was able to add an option with no key. (should not be possible)")
	}
}

func TestAddOption_add_empty_value(t *testing.T) {
	prmpt := NewPrompt("TestQuestion")

	added := prmpt.AddOption("1", "")
	if added || len(prmpt.Blunders.Reported) < 1 {
		t.Errorf("Was able to add an option with no value. (should not be possible)")
	}
}

//////////////////////////////////////////////////////////////////
// User Prompting Tests
//////////////////////////////////////////////////////////////////

func QuickPromptTest(t *testing.T) {
	//need to figure out a way to simulate os.Stdin
}