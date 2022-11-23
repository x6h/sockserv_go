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
    // maximum connections/clients allowed to be connected
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
        deny_connection := false

        // accept incoming connection
        connection, connection_error := listener.Accept()

        for i := 0; i < SERVER_MAX_CONNECTIONS; i++ {
            // if there are no slots available, deny further connection handling
            if i == SERVER_MAX_CONNECTIONS - 1 && connection_list[i] != nil {
                deny_connection = true
            }
        }

        if deny_connection {
            connection.Close()
            continue
        }

        if connection_error != nil {
            fmt.Printf("failed to accept incoming connection. (error: %s)\n", connection_error)
        }

        // start thread to handle the newly accepted connection
        go handle_connection(connection)
    }
}

func handle_connection(connection net.Conn) {
    // add client to the client list
    for i := 0; i < SERVER_MAX_CONNECTIONS; i++ {
        if connection_list[i] == nil {
            connection_list[i] = connection
            active_connections++
            break
        }
    }

    fmt.Printf("server (con): active: %d\n", active_connections)

    for {
        // read message from the client
        message := make([]byte, SERVER_MSG_LENGTH)
        _, read_error := connection.Read(message)

        // send client message to other clients
        for i := 0; i < SERVER_MAX_CONNECTIONS; i++ {
            // skip invalid connections
            if (connection_list[i] == nil) {
                continue
            }

            // skip the author of the message
            if connection_list[i] == connection {
                continue
            }

            _, write_error := connection_list[i].Write(message)

            if write_error != nil {
                fmt.Printf("failed to send message to clients. (error: %s)\n", write_error)
            }
        }

        if read_error == nil {
            fmt.Printf("server (msg): %s\n", string(message))
        } else {
            handle_disconnection(connection)
            break
        }
    }
}

func handle_disconnection(connection net.Conn) {
    for i := 0; i < SERVER_MAX_CONNECTIONS; i++ {
        if connection_list[i] == connection {
            // close connection, not sure if this affects anything tho
            connection.Close()
            // remove connection from connection list
            connection_list[i] = nil
            active_connections--
        }
    }

    fmt.Printf("server (dis): active: %d\n", active_connections)
}
