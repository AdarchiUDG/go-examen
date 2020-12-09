package server

import (
	"errors"
	"fmt"
)

type ChatMessage struct {
	Nick string
	Content string
}

type Chat struct {
	Messages map[string][]string
}

func (c *Chat) Register(nick string, reply *bool) error {
	if c.Messages[nick] != nil {
		*reply = false
		return errors.New("Ya existe un usuario con ese nick")
	}

	*reply = true
	fmt.Printf("Se conecto %s\n", nick)
	return nil
}

func (c *Chat) SendMessage(msg *ChatMessage, reply *bool) error {
	var content string
	content +=  msg.Nick + ": " + msg.Content

	fmt.Println(msg.Nick, msg.Content)

	for user, messages := range c.Messages	{
		c.Messages[user] = append(messages, content)
	}

	*reply = true

	return nil
}

func (c *Chat) CheckMessages(nick string, reply *[]string) error {
	*reply = c.Messages[nick]
	c.Messages[nick] = nil
	
	return nil
}