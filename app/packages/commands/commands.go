package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/packages/dredis"
	"github.com/codecrafters-io/redis-starter-go/app/packages/resp"
)

var dRedis = make(map[string]dredis.CacheItem)

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
		res, err := handleSet(args)
		return res, err
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
	if len(args) < 2 || len(args)%2 != 0 {
		return nil, errors.New("ERR wrong number of arguments for 'set' command")
	}
	key := args[0]
	value := args[1]
	cacheItem := dredis.CacheItem{
		Value:     value,
		CacheType: "string",
		ExpiresAt: -1,
	}

	var err error
	for i := 2; i < len(args); i += 2 {
		arg := args[i]
		argValue := args[i+1]

		cacheItem, err = setArguments(cacheItem, arg, argValue)
		if err != nil {
			return resp.EncodeError("Invalid arguments for 'SET'"), err
		}
	}

	dRedis[key] = cacheItem
	return resp.Encode(rESPONSE_OK), nil
}

func setArguments(cacheItem dredis.CacheItem, arg string, argValue string) (dredis.CacheItem, error) {
	arg = strings.ToUpper(arg)
	switch arg {
	case "PX":
		ttl, err := strconv.Atoi(argValue)
		if err != nil {
			return cacheItem, errors.New("ERR invalid value " + argValue + " for expiry")
		}
		if ttl < -1 {
			return cacheItem, errors.New("ERR invalid expire time in SET")
		}
		if ttl == -1 {
			cacheItem.ExpiresAt = -1
		} else {
			cacheItem.ExpiresAt = time.Duration(time.Now().UnixMilli() + int64(ttl))
		}
	}
	return cacheItem, nil
}

func handleGet(args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("ERR wrong number of arguments for 'get' command")
	}

	value, ok := dRedis[args[0]]
	if !ok || (value.ExpiresAt != -1 && value.ExpiresAt < time.Duration(time.Now().UnixMilli())) {
		return resp.EncodeBulkString(rESPONSE_EMPTY), nil
	}

	return resp.Encode(value.Value), nil
}
