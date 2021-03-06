package utils

import (
	"fmt"
	"sort"

	"github.com/gamedb/gamedb/pkg/chatbot"
)

type chatCommands struct{}

func (chatCommands) name() string {
	return "output-chat-commands"
}

func (chatCommands) run() {

	var commands = chatbot.CommandRegister

	sort.Slice(commands, func(i, j int) bool {
		if commands[i].Type().Order() == commands[j].Type().Order() {
			return commands[i].ID() < commands[j].ID()
		}
		return commands[i].Type().Order() < commands[j].Type().Order()
	})

	lastType := chatbot.CommandType("")
	for _, v := range commands {

		if lastType != v.Type() {
			fmt.Println()
			fmt.Println("&nbsp;")
			fmt.Println()
			fmt.Println(v.Type() + " Commands | Description")
			fmt.Println("--- | ---")
		}

		fmt.Println("/" + v.ID() + " | " + v.Description())

		lastType = v.Type()
	}

	fmt.Println()

	imgurIDs := []string{"g6oINcF", "Fc2ANPn", "dqolEon", "JhQEpJE", "5tS66Aw", "Su20AZ6", "xJx34ZW", "bT0Pcdv", "g16Uekz", "mWoDw3F"}

	for _, id := range imgurIDs {
		fmt.Println("![Global Steam Discord Bot](https://i.imgur.com/" + id + ".png)<br>")
	}
}
