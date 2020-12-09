package main

import (
	"net"
	"net/rpc"
	"fmt"
	"os"
	"./server"
)

func main() {
	chat := new(server.Chat)
	chat.Messages = make(map[string][]string)
	
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
	
  for {
			conn, err := listener.Accept()
      if err != nil {
          continue
			}
      go rpc.ServeConn(conn)
  }
}