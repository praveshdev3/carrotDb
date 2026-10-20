package server

import (
	"io"
	"log"
	"net"

	"github.com/praveshdev3/carrotDb/config"
)

func RunSyncTCPServer() {
	log.Println("Starting the carrot 🥕 server on " + config.Host + ":" + config.Port)

	con_clients := 0

	lsnr, err := net.Listen("tcp", config.Host+":"+config.Port)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}

	for {

		conn, err := lsnr.Accept()
		if err != nil {
			log.Fatal("Error accepting the connection: ", err)
		}

		con_clients++
		log.Println("client connected with address:", conn.RemoteAddr(), "concurrent clients", con_clients)

		for {
			cmd, err := readCommand(conn)
			if err != nil {
				conn.Close()
				con_clients--
				log.Println("client disconnected with address:", conn.RemoteAddr(), "concurrent clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			log.Println("COMMAND: ", cmd)
			if err = respond(cmd, conn); err != nil {
				log.Print("err write:", err)
			}
		}
	}

}

func readCommand(conn net.Conn) (string, error) {
	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection: ", err)
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(cmd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		log.Println("Error writing to connection: ", err)
		return err
	}
	return nil
}
