package main

import (
    "fmt"
    "net"
    "os"
)

// server settings
const (
    SERVER_TYPE = "tcp"
    SERVER_HOST = "localhost"
    SERVER_PORT = "8080"
    // maximum recieve message length
    // (overflow has been observed to be read on the next "connection.Read" call)
    SERVER_MSG_LENGTH = 512
)

func main() {
    // listen for incoming connections
    listener, listener_error := net.Listen(SERVER_TYPE, SERVER_HOST + ":" + SERVER_PORT)

    if listener_error != nil {
        fmt.Printf("failed to listen for incoming connections. (error: %s)\n", listener_error)
        os.Exit(1)
    }

    for {
        // accept incoming connection
        connection, connection_error := listener.Accept()

        if connection_error != nil {
            fmt.Printf("failed to accept incoming connection. (error: %s)\n", connection_error)
        }

        // start goroutine to handle accepted connection while this loop can keep accepting new connections
        go handleConnection(connection)
    }
}

func handleConnection(connection net.Conn) {
    for {
        recieved_data := make([]byte, SERVER_MSG_LENGTH)
        _, read_error := connection.Read(recieved_data)

        if read_error == nil {
            fmt.Println(string(recieved_data))
        } else {
            fmt.Printf("connection error. (error: %s)\n", read_error)
            break
        }
    }
}
