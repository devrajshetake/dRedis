package commands

import (
	"errors"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/packages/resp"
)

var dRedis = make(map[string]string)

func Execute(command string, args []string) ([]byte, error) {
	command = strings.ToUpper(command)
	switch command {
	case "PING":
		res, _ := handlePing()
		return res, nil
	case "ECHO":
		res, _ := handleEcho(args)
		return res, nil
	case "SET":
		res, _ := handleSet(args)
		return res, nil
	case "GET":
		res, _ := handleGet(args)
		return res, nil
	default:
		return nil, errors.New("ERR unknown command '" + command + "'")
	}
}

func handlePing() ([]byte, error) {
	pong := resp.Encode(rESPONSE_PONG)
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

func handleSet(args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("ERR wrong number of arguments for 'set' command")
	}

	dRedis[args[0]] = args[1]
	return resp.Encode(rESPONSE_OK), nil
}

func handleGet(args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("ERR wrong number of arguments for 'get' command")
	}

	value, ok := dRedis[args[0]]
	if !ok {
		return resp.EncodeBulkString(rESPONSE_EMPTY), nil
	}
	return resp.Encode(value), nil
}
