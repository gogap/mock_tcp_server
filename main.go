package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/gogap/mock_tcp_server/config"
)

var (
	conf = config.MockServerConfig{}
)

func main() {

	if bFile, e := ioutil.ReadFile("server.conf"); e != nil {
		log.Fatalln(e)
		return
	} else {
		if e := json.Unmarshal(bFile, &conf); e != nil {
			log.Fatalln(e)
			return
		}
	}

	var tcpAddr *net.TCPAddr
	if addr, e := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", conf.Host, conf.Port)); e != nil {
		log.Fatalln(e)
		return
	} else {
		tcpAddr = addr
	}

	var tcpListener *net.TCPListener
	if listener, e := net.ListenTCP("tcp", tcpAddr); e != nil {
		log.Fatalln(e)
		return
	} else {
		tcpListener = listener
	}

	fmt.Printf("Listening:%s:%d\n", conf.Host, conf.Port)

	for {
		if conn, e := tcpListener.Accept(); e != nil {
			log.Fatalln(e)
			continue
		} else {
			handle_client(conn)
			conn.Close()
		}
	}
}

func handle_client(conn net.Conn) {
	fmt.Printf("client connected: %s\n", conn.RemoteAddr().String())

	var buf [2048]byte

	if _, e := conn.Read(buf[0:]); e != nil {
		log.Fatalln(e)
		return
	} else {
		for _, matchItem := range conf.Matches {
			matched := false
			if matchItem.Type == "string" {
				if strings.Contains(string(buf[0:]), matchItem.MatchData) {
					matched = true
				}
			} else if matchItem.Type == "byte" {
				if bMath, e := hex.DecodeString(matchItem.MatchData); e != nil {
					log.Fatalln(e)
					return
				} else if bytes.Contains(buf[0:], bMath) {
					matched = true
				}
			} else {
				log.Fatalln("unknown match item type, only could be [string|byte]")
				return
			}

			if matched {
				fmt.Printf("[matched %s:%s]\n%s\n", matchItem.Type, matchItem.MatchData, matchItem.ResponseFile)
				if fd, e := ioutil.ReadFile(matchItem.ResponseFile); e != nil {
					log.Fatalln(e)
					return
				} else {
					if _, e := conn.Write(fd); e != nil {
						log.Fatalln(e)
						return
					}
					return
				}
			}
		}
		fmt.Println("nothing matched.")
	}
}
