package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/genesixx/coalibot/Miscs"
	"github.com/genesixx/coalibot/Struct"

	"github.com/nlopes/slack"
)

var commands = map[string]func(string, *Struct.Message) bool{
	"hello":    Miscs.Hello,
	"vdm":      Miscs.Vdm,
	"roulette": Miscs.Roulette,
	"coin":     Miscs.Coin,
	"meteo":    Miscs.Meteo,
	"roll":     Miscs.Roll,
}

func handleCommand(event *Struct.Message) {
	var isCommand = false
	var option = ""
	var command = ""

	event.Message = strings.Join(strings.Fields(event.Message), " ")
	fmt.Printf("<#%s> @%s: %s\n", event.Channel, event.User, event.Message)
	splited := strings.Split(event.Message, " ")
	if indexOf(splited[0], []string{"coalibot", "bc", "cb"}) > -1 && len(splited) > 1 {
		command = splited[1]
		option = strings.Join(splited[2:], " ")
		isCommand = reply(command, event)
		if !isCommand && commands[command] != nil {
			isCommand = commands[strings.ToLower(command)](option, event)
		}
	} else if splited[0][0] == '!' && len(splited[0]) > 1 {
		command = splited[0][1:]
		option = strings.Join(splited[1:], " ")
		isCommand = reply(command, event)
		if !isCommand && commands[command] != nil {
			isCommand = commands[strings.ToLower(command)](option, event)
		}
	}
	fmt.Printf("command %s option %s\n", command, option)
}

func reply(command string, event *Struct.Message) bool {
	// Open our jsonFile
	jsonFile, err := os.Open("reply.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// a map container to decode the JSON structure into
	c := make(map[string]interface{})

	// unmarschal JSON
	e := json.Unmarshal(byteValue, &c)
	if e != nil || c[command] == nil {
		return false
	}

	// output result to STDOUT
	fmt.Printf("reply %s\n", c[command].(string))
	event.API.PostMessage(event.Channel, c[command].(string), slack.PostMessageParameters{})
	return true
}

func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}