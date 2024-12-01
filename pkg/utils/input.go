package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
