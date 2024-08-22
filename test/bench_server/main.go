package main

import (
	"github.com/Eresh-tech/smgn"
	"github.com/Eresh-tech/smgn/codec"
	"time"
)

var log = smgn.GetLogger()

type Handler struct{}

func (*Handler) OnConnect(c *smgn.Conn) {
	log.Info("connect:", c.GetFd(), c.GetAddr())
}
func (*Handler) OnMessage(c *smgn.Conn, bytes []byte) {
	c.WriteWithEncoder(bytes)
	log.Info("read:", string(bytes))
}
func (*Handler) OnClose(c *smgn.Conn, err error) {
	log.Info("close:", c.GetFd(), err)
}

func main() {
	server, err := smgn.NewServer(":8080", &Handler{},
		smgn.WithDecoder(codec.NewHeaderLenDecoder(2)),
		smgn.WithEncoder(codec.NewHeaderLenEncoder(2, 1024)),
		smgn.WithTimeout(5*time.Second),
		smgn.WithReadBufferLen(10))
	if err != nil {
		log.Info("err")
		return
	}

	server.Run()
}
