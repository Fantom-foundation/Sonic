package prompt

import geth_prompt "github.com/ethereum/go-ethereum/console/prompt"

//go:generate mockgen -source prompt.go -destination prompt_mock.go -package prompt

// UserPrompt is the default prompter used by the console to prompt the user for
// various types of inputs.
// In enables tests to replace the prompter with a mock.
var UserPrompt UserPrompter = geth_prompt.Stdin

// UserPrompter is a re-export of the geth_prompt.UserPrompter interface.
// It is used to generate mocks for tests.
type UserPrompter interface {
	PromptInput(prompt string) (string, error)
	PromptPassword(prompt string) (string, error)
	PromptConfirm(prompt string) (bool, error)
	SetHistory(history []string)
	AppendHistory(command string)
	ClearHistory()
	SetWordCompleter(completer geth_prompt.WordCompleter)
}

// static assert: user prompter declared in this file must implement the one in geth_prompt
var _ geth_prompt.UserPrompter = UserPrompt
