package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

// AskYesNo prompts the user with a yes/no question and returns the boolean response
func AskYesNo(question string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", question)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

// AskInput prompts the user to enter a text input
func AskInput(question string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + ": ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

// SelectOption presents a list of options for the user to choose from
func SelectOption(question string, options []string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	for i, option := range options {
		fmt.Printf("%d: %s\n", i+1, option)
	}
	fmt.Print("Enter the number of your choice: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || choice < 1 || choice > len(options) {
		return "", fmt.Errorf("invalid choice")
	}
	return options[choice-1], nil
}

// MultiSelectOptions allows the user to select multiple options from a list
func MultiSelectOptions(question string, options []string) ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	for i, option := range options {
		fmt.Printf("%d: %s\n", i+1, option)
	}
	fmt.Print("Enter the numbers of your choices (comma-separated): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	choices := strings.Split(strings.TrimSpace(input), ",")
	var selected []string
	for _, choice := range choices {
		num, err := strconv.Atoi(strings.TrimSpace(choice))
		if err != nil || num < 1 || num > len(options) {
			return nil, fmt.Errorf("invalid choice")
		}
		selected = append(selected, options[num-1])
	}
	return selected, nil
}

// AskPassword prompts for a password input with hidden characters
func AskPassword(question string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + ": ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}
