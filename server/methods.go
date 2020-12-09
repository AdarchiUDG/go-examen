package server

import (
	"errors"
	"fmt"
	"github.com/jroimartin/gocui"
)

type ChatMessage struct {
	Nick string
	Content string
}

type ChatFile struct {
	Bytes []byte
	Name string
	Nick string
}

type Chat struct {
	AllMessages []string
	Messages map[string][]string
	Files map[string][]byte
	Gui *gocui.Gui
}

func (c *Chat) Register(nick string, reply *bool) error {
	*reply = true

	c.Gui.Update(func(g *gocui.Gui) error {
		chatView, _ := g.View("chat")
		fmt.Fprintln(chatView, "Se conecto", nick)
		return nil
	})
	return nil
}

func (c *Chat) SendMessage(msg *ChatMessage, reply *bool) error {
	var content string
	content +=  msg.Nick + ": " + msg.Content

	for user, messages := range c.Messages	{
		c.Messages[user] = append(messages, content)
	}

	c.AllMessages = append(c.AllMessages, content)
	c.Gui.Update(func(g *gocui.Gui) error {
		chatView, _ := g.View("chat")
		fmt.Fprintln(chatView, content)
		return nil
	})
	*reply = true

	return nil
}

func (c *Chat) SendFile(msg *ChatFile, reply *bool) error {
	content := fmt.Sprintf("* %s envio un archivo: %s\n* Usa /dl [nombre] para descargarlo", msg.Nick, msg.Name)
	for user, messages := range c.Messages	{
		c.Messages[user] = append(messages, content)
	}

	c.Files[msg.Name] = msg.Bytes
	c.AllMessages = append(c.AllMessages, content)

	c.Gui.Update(func(g *gocui.Gui) error {
		chatView, _ := g.View("chat")
		fmt.Fprintln(chatView, content)
		return nil
	})

	*reply = true

	return nil
}

func (c *Chat) GetFile(name string, reply *[]byte) error {
	if val, ok := c.Files[name]; ok  {
		*reply = val
	
		return nil
	}

	return errors.New("No se encontro el archivo")
}

func (c *Chat) GetFileNames(a bool, reply *[]string) error {
	var names [] string

	for k, _ := range(c.Files) {
		names = append(names, k)
	}

	if len(names) > 0 {
		*reply = names
		return nil
	}
	
	return errors.New("!! No hay archivos")
}

func (c *Chat) CheckMessages(nick string, reply *[]string) error {
	*reply = c.Messages[nick]
	c.Messages[nick] = nil
	
	return nil
}