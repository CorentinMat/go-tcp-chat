package main

import (
	"log"
	"net"
)

func main() {
	server := newServer()
	go server.run()
	listenner, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalln("listenner error")
	}

	defer listenner.Close()
	log.Printf("started listening on 8888")
	for {
		conn, err := listenner.Accept()
		if err != nil {
			log.Fatal("connection error : ", err)
			continue
		}
		go server.newClient(conn)
	}
}
