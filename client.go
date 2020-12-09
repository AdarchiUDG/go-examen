package main

import (
	"bufio"
	"os"
	"fmt"
	"net/rpc"
	"log"
	"strings"
	"time"
)

type ChatMessage struct {
	Nick string
	Content string
}

func main() {
	var allMessages []string
	option := 0
	in := bufio.NewReader(os.Stdin)
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

	go func () {
		var messages []string

		for {
			err := client.Call("Chat.CheckMessages", nickname, &messages)
			if err != nil {
				log.Fatalln("El servidor se ha cerrado.")
			}

			for _, v := range messages {
				fmt.Println(v)
				allMessages = append(allMessages, v)
			}

			time.Sleep(time.Second)
		}
	}()

	for option != 4{
		fmt.Println("1. Enviar Mensaje")
		fmt.Println("2. Enviar Archivo")
		fmt.Println("3. Mostrar Chat")
		fmt.Println("4. Salir")
		fmt.Println("Teclea el numero de la opcion deseada")
	
		fmt.Scanf("%d", &option)
		fmt.Scanln()

		switch option {
			case 1:
				var reply bool

				message := ChatMessage{
					Nick: nickname,
					Content: readLine(in, "Mensaje: ") }

				client.Call("Chat.SendMessage", message, &reply)
			case 2:
				
			case 3:
				for _, v := range(allMessages) {
					fmt.Println(v)
				}
			default:
		}

		if option != 4 {
			fmt.Scanln()
		}
	}
	fmt.Println("Saliendo . . .")
}

func readLine(in *bufio.Reader, str string) string {
	fmt.Print(str)
	line, _ := in.ReadString('\n')
	return strings.TrimSuffix(line, "\r\n")
}