package main

import (
	"bufio"
	"os"
	"fmt"
	"net/rpc"
	"log"
	"strings"
	"time"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"path"
	"./gui"
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

func main() {
	var allMessages []string
	in := bufio.NewReader(os.Stdin)

	/* Starting RPC connection */
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	defer client.Close()

	nickname := readLine(in, "Nickname: ")
	var reply bool
	fmt.Println("Conectando . . .")
	err = client.Call("Chat.Register", nickname, &reply)

	if err != nil {
		log.Fatal(err)
	}
	
	/* Starting GUI */
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(gui.Layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, func (g *gocui.Gui, view *gocui.View) error {
		var reply bool
		in, _ := g.View("input")
		chat, _ := g.View("chat")
		message := strings.TrimSuffix(in.ViewBuffer(), "\n")
		
		if strings.HasPrefix(message, "/exit") {
			os.Exit(1)
		} else if strings.HasPrefix(message, "/file ") {
			message = strings.TrimPrefix(message, "/file ")
			file, err := os.Open(message)
			if err != nil {
				fmt.Fprintln(chat, "!! Error: El archivo \"" + message + "\" no se encontro.")
			} else {
				bytes, _ := ioutil.ReadAll(file)
				client.Call("Chat.SendFile", ChatFile{
					Nick: nickname,
					Bytes: bytes,
					Name: message,
				}, &reply)
			}
		} else if (strings.HasPrefix(message, "/dl ")) {
			var file []byte
			message = strings.TrimPrefix(message, "/dl ")
			err := client.Call("Chat.GetFile", message, &file)
			if err != nil {
				fmt.Fprintln(chat, "!! Error: El archivo \"" + message + "\" no se encontro.")
			} else {
				userFolder := "./downloads/" + nickname
				if _, err := os.Stat(userFolder); os.IsNotExist(err) {
					os.MkdirAll(userFolder, 775)
				}

				p := path.Join(userFolder, message)
				err := ioutil.WriteFile(p, file, 0664)
				if err != nil {
					fmt.Fprintln(chat, err)
				} else {
					fmt.Fprintln(chat, "!! Se guardo el archivo exitosamente en " + p)
				}
			}
		
		} else if (strings.HasPrefix(message, "/files")) {
			var files []string
			err := client.Call("Chat.GetFileNames", true, &files)
			if err != nil {
				fmt.Fprintln(chat, err)
			} else {
				fmt.Fprintln(chat, "--- Archivos en el servidor --")
				for _, file := range(files) {
					fmt.Fprintln(chat, file)
				}
			}
		} else {
			client.Call("Chat.SendMessage", ChatMessage{
				Nick: nickname,
				Content: message,
			}, &reply)
		}
		
	
		in.Clear()
		in.SetCursor(0, 0)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	go func () {
		var messages []string

		for {
			err := client.Call("Chat.CheckMessages", nickname, &messages)
			if err != nil {
				log.Fatalln("El servidor se ha cerrado.")
			}

			for _, v := range messages {
				allMessages = append(allMessages, v)
				g.Update(func (g *gocui.Gui) error {
					view, _ := g.View("chat")
					fmt.Fprintln(view, v)
					return nil
				})
			}

			time.Sleep(time.Second)
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func readLine(in *bufio.Reader, str string) string {
	fmt.Print(str)
	line, _ := in.ReadString('\n')
	return strings.TrimSuffix(line, "\r\n")
}