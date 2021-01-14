package server

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

// RunServer is ..
func RunServer(address string) {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		fmt.Println(err)
		return
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	l, err := tls.Listen("tcp", address, &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	ConnectionContainer := &sync.Map{}
	UserContainer := make(map[string]string)
	ChatContainer := make(map[string]string)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		remoteAddress := c.RemoteAddr().String()
		tlscon, ok := c.(*tls.Conn)
		if ok {
			log.Print("ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		fmt.Printf("Serving %s\n", remoteAddress)
		c.Write([]byte("Per cominciare dimmi il tuo nome per favore!\n"))
		go func() {
			for {
				netData, readErr := bufio.NewReader(c).ReadString('\n')
				temp := strings.TrimSpace(string(netData))
				if readErr != nil {
					fmt.Println(err)
				}
				remoteAddress := c.RemoteAddr().String()
				fmt.Printf("Receiving from: %s\n", remoteAddress)
				user, usernameSetted := UserContainer[remoteAddress]
				if !usernameSetted {
					fmt.Println("setting username: ", temp)
					response := fmt.Sprintf("Ciao %s, Con chi vuoi parlare?\n", temp)
					_, connFound := ConnectionContainer.Load(temp)
					if connFound {
						c.Write([]byte("Il nome è già in uso provane un'altro!\n"))
					} else {
						UserContainer[remoteAddress] = temp
						ConnectionContainer.Store(temp, c)
						c.Write([]byte(response))
					}
				} else {
					if readErr == io.EOF {
						fmt.Println("deleting user connections: ", user)
						ConnectionContainer.Delete(user)
						delete(UserContainer, remoteAddress)
						return
					}
					fmt.Println("ricerca chat: ", user)
					remoteChat, chatFound := ChatContainer[user]
					fmt.Println("trovato: ", remoteChat)
					if !chatFound {
						fmt.Println("chat setting: ", temp)
						ChatContainer[user] = temp
						c.Write([]byte(string("Chat avviata:\n")))
					} else {
						fmt.Println("tentativo invio: ", remoteChat)
						conn, ok := ConnectionContainer.Load(remoteChat)
						fmt.Println("risultato: ", ok)
						if ok {
							message := fmt.Sprintf("message from %s: %s\n", user, temp)
							fmt.Println(message)
							conn.(net.Conn).Write([]byte(message))
						} else {
							c.Write([]byte(string("Utente non disponibile!:\n")))

						}
					}
				}
				fmt.Println("received:", temp)
			}
		}()
	}
}
