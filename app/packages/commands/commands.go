package commands

import (
	"errors"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/packages/resp"
)

func Execute(command string, args []string) ([]byte, error) {
	command = strings.ToUpper(command)
	switch command {
	case "PING":
		res, _ := handlePing()
		return res, nil
	case "ECHO":
		res, _ := handleEcho(args)
		return res, nil
	default:
		return nil, errors.New("ERR unknown command '" + command + "'")
	}
}

func handlePing() ([]byte, error) {
	pong := resp.Encode("PONG")
	if pong == nil {
		return nil, errors.New("ERR failed to encode response")
	}
	return pong, nil
}

func handleEcho(args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("ERR wrong number of arguments for 'echo' command")
	}

	echo := resp.Encode(args[0])
	if echo == nil {
		return nil, errors.New("ERR failed to encode response")
	}
	return echo, nil
}
