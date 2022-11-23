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
    // maximum connections
    SERVER_MAX_CONNECTIONS = 3
    // maximum recieve message length
    // (overflow has been observed to be read on the next "connection.Read" call)
    SERVER_MSG_LENGTH = 512
)

// globals
var connection_list [SERVER_MAX_CONNECTIONS]net.Conn
var active_connections int = 0

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
            continue
        }

        // don't continue if max connections have been reached
        if active_connections == SERVER_MAX_CONNECTIONS {
            connection.Close()
            continue
        }

        // start thread to handle the newly accepted connection
        go HandleConnection(connection)
    }
}
