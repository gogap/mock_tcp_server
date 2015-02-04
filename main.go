package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gogap/mock_tcp_server/config"
)

var (
	conf = config.MockServerConfig{}

	reqID   int64 = 0
	dumpDir       = "./dump"

	configFile = "server.conf"
)

func main() {

	log.SetOutput(os.Stdout)

	short := " (shorthand)"

	configfileUsage := "the server config file"
	flag.StringVar(&configFile, "config", "", configfileUsage)
	flag.StringVar(&configFile, "c", "", configfileUsage+short)

	flag.Parse()

	now := time.Now()

	if configFile == "" {
		configFile = "server.conf"
	}

	if bFile, e := ioutil.ReadFile(configFile); e != nil {
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

	if conf.DumpRequest {
		dumpDir = fmt.Sprintf("./dump/%d", now.UnixNano())
		fmt.Printf("[dump dir]: %s\n", dumpDir)
		if !is_dir_exist(dumpDir) {
			if e := os.MkdirAll(dumpDir, os.ModePerm); e != nil {
				log.Fatal(e)
				return
			}
		}
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
	atomic.AddInt64(&reqID, 1)

	fmt.Printf("[%d]client connected: %s\n", reqID, conn.RemoteAddr().String())

	var buf [2048]byte

	if l, e := conn.Read(buf[0:]); e != nil {
		log.Fatalln(e)
		return
	} else {
		if conf.DumpRequest {
			filename := fmt.Sprintf("%s/%d.dat", dumpDir, reqID)
			ioutil.WriteFile(filename, buf[0:l], 0666)
		}
		for _, matchItem := range conf.Matches {
			matched := false
			if matchItem.Type == "string" {
				if strings.Contains(string(buf[0:l]), matchItem.MatchData) {
					matched = true
				}
			} else if matchItem.Type == "byte" {
				if bMath, e := hex.DecodeString(matchItem.MatchData); e != nil {
					log.Fatalln(e)
					return
				} else if bytes.Contains(buf[0:l], bMath) {
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

func is_dir_exist(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}
