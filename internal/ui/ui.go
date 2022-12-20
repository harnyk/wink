package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type UI interface {
	AskString(prompt string) (string, error)
	AskPassword(prompt string) (string, error)
}

func NewUI() UI {
	return &ui{}
}

type ui struct {
}

func (u *ui) AskString(prompt string) (string, error) {
	fmt.Println(prompt)
	var result string
	fmt.Scanln(&result)
	return result, nil
}

func (u *ui) AskPassword(prompt string) (string, error) {
	fmt.Println(prompt)

	password, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		return "", err
	}

	return string(password), nil
}
