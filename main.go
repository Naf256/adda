package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"
)

const (
	Port = "6969"
	BanLimit = 10*60.0
	SafeMode = false
	MessageRate = 1.0
	StrikeLimit = 10
)

type MessageType int

const (
	ClientConnected MessageType = iota + 1
	ClientDisconnected
	NewMessage
)


type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

type Client struct {
	Conn net.Conn
	LastMessage time.Time
	Strikes int
}

// help printout ip while debugging
func sensitive(ip string) string {
	if SafeMode {
		return "[PRIVATE_IP]"
	}

	return ip
}

func client(conn net.Conn, messages chan Message) {
	buffer := make([]byte, 512)

	for {
		// keep receiving any message sent from client
		n, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			messages <- Message{
				Type: ClientDisconnected,
				Conn: conn,
			}
			return
		}

		text := string(buffer[0:n])
		messages <- Message{
			Type: NewMessage,
			Conn: conn,
			Text: text,
		}
	}
}

func server(messages chan Message) {
	clients := map[string]*Client{}
	bannedBots := map[string]time.Time{}

	for {
		msg := <- messages
		switch msg.Type {

		case ClientConnected:
			addr := msg.Conn.RemoteAddr().String()
			now := time.Now()
			bannedAt, banned := bannedBots[addr]
			if banned {
				if now.Sub(bannedAt).Seconds() >= BanLimit {
					delete(bannedBots, addr)
					banned = false
				}
			}
			
			if !banned {
				log.Printf("client: %s connected", sensitive(addr))
				clients[addr] = &Client{
					Conn: msg.Conn,
					LastMessage: time.Now(),
				}
			} else {
				msg.Conn.Write([]byte(fmt.Sprintf("You are banned Bot: %f secs left\n", BanLimit - float64(now.Sub(bannedAt).Seconds()))))
				msg.Conn.Close()
			}

		case ClientDisconnected:
			addr := msg.Conn.RemoteAddr().String()
			log.Printf("Client %s disconnected", sensitive(addr))
			delete(clients, addr)

		case NewMessage:
			addr := msg.Conn.RemoteAddr().String()
			now := time.Now()
			author := clients[addr]
			if author != nil {
				// check if client is sending too many messages if then bann him
				if now.Sub(author.LastMessage).Seconds() >= MessageRate {
					// check if message sent by user is valid is valid
					if utf8.ValidString(msg.Text) {
						author.LastMessage = now
						author.Strikes = 0
						log.Printf("client %s sent message %s", sensitive(addr), msg.Text)
						for _, client := range clients {
							if client.Conn.RemoteAddr().String() != addr {
								client.Conn.Write([]byte(msg.Text))
							}
						}
					} else {
						author.Strikes += 1
						if author.Strikes >= StrikeLimit {
							bannedBots[addr] = time.Now()
							author.Conn.Write([]byte("You are banned Bot\n"))
							author.Conn.Close()
						}
					}
				} else {
					author.Strikes += 1
					if author.Strikes >= StrikeLimit {
						bannedBots[addr] = time.Now()
						author.Conn.Write([]byte("You are banned Bot\n"))
						author.Conn.Close()
					}
				}
			} else {
				author.Conn.Close()
			}
		}

	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalln("Failed to start a tcp connection...")
	}


	log.Printf("Listening to TCP connections on port %s ...\n", Port);

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting a connection")
			continue
		}

		messages <- Message{
			Type: ClientConnected,
			Conn: conn,
		}

		go client(conn, messages)
	}
}
