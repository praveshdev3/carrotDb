package core

import (
	"fmt"
	"net"
)

func Respond(cmd *RedisCmd, conn net.Conn) error {
	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd, conn)
	default:
		return fmt.Errorf("unknown command: %s", cmd.Cmd)
	}
}

func evalPing(cmd *RedisCmd, conn net.Conn) error {
	var b []byte

	if len(cmd.Args) >= 2 {
		return fmt.Errorf("ERR wrong number of arguments for 'ping' command")
	}

	if len(cmd.Args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(cmd.Args[0], false)
	}

	conn.Write(b)
	return nil
}
