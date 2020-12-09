package main

import (
	"net"
	"net/rpc"
	"fmt"
	"os"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
	"io/ioutil"
	"path"
	"./server"
	"./gui"
)

func main() {
	chat := new(server.Chat)
	chat.Messages = make(map[string][]string)
	chat.Files = make(map[string][]byte)
	
  rpc.Register(chat)

  tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
  if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}	

  listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}

	/* Starting GUI */
	g, err := gocui.NewGui(gocui.OutputNormal)
	chat.Gui = g
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(gui.Layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, func (g *gocui.Gui, view *gocui.View) error {
		in, _ := g.View("input")
		c, _ := g.View("chat")
		message := strings.TrimSuffix(in.ViewBuffer(), "\n")
		
		switch message {
			case "exit":
				fmt.Fprintln(c, "Cerrando el servidor . . .\n")
				os.Exit(1)
			case "files":
				for n, _ := range(chat.Files) {
					fmt.Fprintln(c, n)
				}
			case "backup":
				fmt.Fprintln(c, "Respaldando mensajes. . .")
				var messages string
				
				for _, m := range(chat.AllMessages) {
					messages += m + "\n"
				}

				ioutil.WriteFile("chat.txt", []byte(messages), 664)

				fmt.Fprintln(c, "Respaldando archivos. . .")
				if _, err := os.Stat("./files"); os.IsNotExist(err) {
					os.Mkdir("./files", 775)
				}

				for n, f := range(chat.Files) {
					p := path.Join("./files", n)
					ioutil.WriteFile(p, f, 664)
				}
			default:
		}
		in.Clear()
		in.SetCursor(0, 0)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	go func () {
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	}()
	
  for {
			conn, err := listener.Accept()
      if err != nil {
          continue
			}
      go rpc.ServeConn(conn)
  }
}