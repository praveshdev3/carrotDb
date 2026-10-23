package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/praveshdev3/carrotDb/config"
	"github.com/praveshdev3/carrotDb/core"
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
			respond(cmd, conn)
		}
	}

}

func readCommand(conn net.Conn) (*core.RedisCmd, error) {
	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection: ", err)
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		log.Println("Error decoding command: ", err)
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  tokens[0],
		Args: tokens[1:],
	}, nil
}

func respond(cmd *core.RedisCmd, conn net.Conn) {
	err := core.Respond(cmd, conn)
	if err != nil {
		respondError(err, conn)
	}
}

func respondError(err error, conn net.Conn) {
	conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}
