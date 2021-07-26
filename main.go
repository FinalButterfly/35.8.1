package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

var proverbs = [19]string{
	"Don't communicate by sharing memory, share memory by communicating.", "Concurrency is not parallelism.",
	"Channels orchestrate; mutexes serialize.",
	"The bigger the interface, the weaker the abstraction.",
	"Make the zero value useful.",
	"interface{} says nothing.",
	"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
	"A little copying is better than a little dependency.",
	"Syscall must always be guarded with build tags.",
	"Cgo must always be guarded with build tags.",
	"Cgo is not Go.",
	"With the unsafe package there are no guarantees.",
	"Clear is better than clever.",
	"Reflection is never clear.",
	"Errors are values.",
	"Don't just check errors, handle them gracefully.",
	"Design the architecture, name the components, document the details.",
	"Documentation is for users.",
	"Don't panic.",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const addr = "localhost:12345"

const proto = "tcp4"

func main() {
	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	handleFunc := func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			handleConn(conn)
		}
	}

	for i := 0; i < 24; i++ {
		go handleFunc()
	}
	wg.Add(1)
	wg.Wait()
}

func handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	b, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := strings.TrimSuffix(string(b), "\n")
	msg = strings.TrimSuffix(msg, "\r")
	if msg == "proverbs" {
		conn.Write([]byte("Here we go!\n"))
	}

	ticker := time.NewTicker(time.Second * 3)

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				n := rand.Intn(18)
				conn.Write([]byte(proverbs[n] + "\n"))
			}
		}
	}()
}
