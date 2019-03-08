package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
)

func main() {
    listener, err := net.Listen("tcp", "localhost:8000")
    if err != nil {
        log.Fatal(err)
    }

    go broadcaster()
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Print(err)
            continue
        }
        go handleConn(conn)
    }
}

type Peer struct {
    Name string
    Channel chan<- string
}

func (p *Peer) Send(message string) {
    p.Channel <- message
}

var (
    entering = make(chan Peer)
    leaving  = make(chan Peer)
    messages = make(chan string)
)

func broadcaster() {
    clients := make(map[Peer]bool)
    for {
        select {
        case msg := <-messages:
            for peer := range clients {
                peer.Send(msg)
            }
        case peer := <-entering:
            if len(clients) == 0 {
                peer.Send("You are the only visitor here!")
            } else {
                peer.Send("Other peers:")
                for other, _ := range clients {
                    peer.Send(other.Name)
                }
            }
            clients[peer] = true
        case peer:= <-leaving:
            delete(clients, peer)
            close(peer.Channel)
        }
    }
}

func handleConn(conn net.Conn) {
    ch := make(chan string)
    go clientWriter(conn, ch)

    who := conn.RemoteAddr().String()
    peer := Peer{Name:who, Channel:ch}

    ch <- "You are " + who
    messages <- who + " has arrived"
    entering <- peer

    input := bufio.NewScanner(conn)
    for input.Scan() {
        messages <- who + ": " + input.Text()
    }

    leaving <- peer
    messages <- who + " has left"
    _ = conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
    for msg := range ch {
        _, _ = fmt.Fprintln(conn, msg)
    }
}