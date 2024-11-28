// internal/commands/commands.go
package commands

import (
	"fmt"
	"os/exec"
	"strings"
)

type CommandVariable interface {
	GetValue() (string, error)
}

type ShellCommandVariable struct {
	Command string
}

func (v *ShellCommandVariable) GetValue() (string, error) {
	cmd := exec.Command("sh", "-c", v.Command)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}
	return string(output), nil
}

type InputCommandVariable struct {
	Message string
}

func (v *InputCommandVariable) GetValue() (string, error) {
	fmt.Print(v.Message + ": ")
	var input string
	_, err := fmt.Scanln(&input)
	return input, err
}

// ProcessPrompt replaces variables in the prompt with their values
func ProcessPrompt(prompt string, variables map[string]CommandVariable) (string, error) {
	result := prompt
	for name, variable := range variables {
		value, err := variable.GetValue()
		if err != nil {
			return "", fmt.Errorf("failed to get value for variable %s: %w", name, err)
		}
		result = strings.ReplaceAll(result, "{{"+name+"}}", value)
	}
	return result, nil
}
