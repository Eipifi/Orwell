package main
import (
    "net"
    "strconv"
    "fmt"
    "os"
)

func serve(port int, handler func(net.Conn)(error)) { // TODO: graceful shutdown
    server, err := net.Listen("tcp", ":" + strconv.Itoa(port))
    if err != nil {
        fmt.Printf("Failed to start server: %s\n", err.Error())
        os.Exit(1)
    }

    defer server.Close()
    fmt.Printf("Listening on port %d\n", port)

    for {
        conn, err := server.Accept()
        if err != nil {
            fmt.Println("Failed to accept connection: %s\n", err.Error())
        } else {
            fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())
            go func(){
                if e := handler(conn); e != nil {
                    fmt.Println("Error: ", e)
                }
                if e := conn.Close(); e != nil {
                    fmt.Println("Failed to close conn: ", e)
                }
            }()
        }
    }
}
