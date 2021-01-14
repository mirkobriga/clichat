package client

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strings"
)

// RunClient is ...
func RunClient(address string) {
	done := make(chan bool)
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	c, err := tls.Dial("tcp", address, &config)
	if err != nil {
		fmt.Println(err)
		done <- true
		return
	}

	state := c.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println(v.Subject)
	}
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	reader := bufio.NewReader(c)
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			fmt.Fprintf(c, text+"\n")
			if strings.TrimSpace(string(text)) == "STOP" {
				done <- true
				return
			}
		}
	}()
	go func() {
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				done <- true
				return
			}
			fmt.Print("->: " + message)
		}
	}()
	<-done
	c.Close()
}
