package main

import (
	"github.com/Eresh-tech/smgn"
	"net"
	"strconv"
	"time"
)

var log = smgn.GetLogger()

type Handler struct{}

func (*Handler) OnConnect(c *smgn.Conn) {
	log.Info("server:connect:", c.GetFd(), c.GetAddr())
}
func (*Handler) OnMessage(c *smgn.Conn, bytes []byte) {
	c.Write(bytes)
	log.Info("server:read:", string(bytes))
}
func (*Handler) OnClose(c *smgn.Conn, err error) {
	log.Info("server:close:", c.GetFd(), err)
}

func startServer() {
	server, err := smgn.NewServer(":8080", &Handler{},
		smgn.WithTimeout(5*time.Second),
		smgn.WithReadBufferLen(10))
	if err != nil {
		log.Info("err")
		return
	}

	server.Run()
}

var space = "                          client:"

func startClient(i int) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Info(space, i, "error dialing", err.Error())
		return // 终止程序
	}

	go func() {
		for {
			buf := make([]byte, 100)
			n, err := conn.Read(buf)
			if err != nil {
				log.Error(space, i, " ", err)
				return
			}
			log.Info(space, i, " ", string(buf[0:n]))
		}
	}()

	for i := 0; i < 10; i++ {
		_, err := conn.Write([]byte("hello" + strconv.Itoa(i)))
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func batchStartClient() {
	for i := 0; i < 1; i++ {
		go startClient(i)
	}
}

func main() {
	go startServer()

	time.Sleep(1 * time.Second)

	go batchStartClient()

	select {}
}
