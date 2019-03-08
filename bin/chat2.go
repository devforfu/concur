package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "time"
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
    Last time.Time
    Conn net.Conn
    Channel chan<- string
}

func (p *Peer) Send(message string) {
    p.Channel <- message
}

func (p *Peer) Close() error {
    if p.Conn != nil {
        err := p.Conn.Close()
        p.Conn = nil
        return err
    }
    return nil
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
            if peer.Conn != nil {
                close(peer.Channel)
            }
        default:
            for peer := range clients {
                if time.Since(peer.Last) > 5 * time.Second {
                    fmt.Printf("Disconnecting idle peer: %s\n", peer.Name)
                    delete(clients, peer)
                    peer.Close()
                }
            }
        }
    }
}

func handleConn(conn net.Conn) {
    ch := make(chan string)
    go clientWriter(conn, ch)

    ch <- "What is your name?"
    input := bufio.NewScanner(conn)
    for {
        if input.Scan() {
            break
        }
    }

    peer := Peer{Name:input.Text(), Channel:ch, Conn:conn, Last:now()}
    messages <- peer.Name + " has arrived"
    entering <- peer

    input = bufio.NewScanner(conn)
    for input.Scan() {
        messages <- peer.Name + ": " + input.Text()
        peer.Last = now()
    }

    leaving <- peer
    messages <- peer.Name + " has left"
    _ = peer.Close()
}

func now() time.Time {
    return time.Now().UTC()
}

func clientWriter(conn net.Conn, ch <-chan string) {
    for msg := range ch {
        _, _ = fmt.Fprintln(conn, msg)
    }
}