package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/packages/commands"
	"github.com/codecrafters-io/redis-starter-go/app/packages/resp"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("dRedis Listening on port 6379...")

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		dataLen, err := conn.Read(buf)

		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}

		// fmt.Printf("Received: %v\n", string(buf[:dataLen]))
		parsed, _ := resp.Parse(buf[:dataLen])
		// fmt.Printf("parsed: %v %v\n", parsed, len(parsed.([]interface{})))

		response := getResponse(parsed)
		conn.Write(response)
	}
}

func getResponse(parsedData interface{}) []byte {
	arr := parsedData.([]interface{})
	if len(arr) == 0 {
		return []byte("-ERR no command provided\r\n")
	}

	command := arr[0].(string)
	args := make([]string, 0)
	// fmt.Printf("command: %v\n args: %v", command, args)
	for i := 1; i < len(arr); i++ {
		args = append(args, arr[i].(string))
	}

	res, err := commands.Execute(command, args)

	if err != nil {
		return []byte("-" + err.Error() + "\r\n")
	}

	return res
}
