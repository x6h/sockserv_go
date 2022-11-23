package main

import (
    "fmt"
    "net"
)

func HandleConnection(connection net.Conn) {
    // add connection to the connection list
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
            HandleDisconnection(connection)
            break
        }
    }
}

func HandleDisconnection(connection net.Conn) {
    for i := 0; i < SERVER_MAX_CONNECTIONS; i++ {
        if connection_list[i] == connection {
            // close connection, not sure if this affects anything tho
            connection_list[i].Close()
            // remove connection from connection list
            connection_list[i] = nil
            active_connections--
        }
    }

    fmt.Printf("server (dis): active: %d\n", active_connections)
}
