package main

import (
	"github.com/Eresh-tech/smgn"
	"github.com/Eresh-tech/smgn/codec"
	"math"
	"net"
	"strconv"
	"time"
)

var (
	decoder = codec.NewHeaderLenDecoder(2)
	encoder = codec.NewHeaderLenEncoder(2, 1024)
)

var log = smgn.GetLogger()

type Handler struct{}

func (*Handler) OnConnect(c *smgn.Conn) {
	log.Info("server:connect:", c.GetFd(), c.GetAddr())
}
func (*Handler) OnMessage(c *smgn.Conn, bytes []byte) {
	c.WriteWithEncoder(bytes)
	log.Info("server:read:", string(bytes))
}
func (*Handler) OnClose(c *smgn.Conn, err error) {
	log.Info("server:close:", c.GetFd(), err)
}

func startServer() {
	server, err := smgn.NewServer(":8080", &Handler{},
		smgn.WithDecoder(decoder),
		smgn.WithEncoder(encoder),
		smgn.WithTimeout(5*time.Second),
		smgn.WithReadBufferLen(20))
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

	buffer := codec.NewBuffer(make([]byte, 1024))

	var handler = func(bytes []byte) {
		log.Info(space, i, " ", string(bytes))
	}

	go func() {
		for {
			_, err := buffer.ReadFromReader(conn)
			if err != nil {
				log.Error(space, i, " ", err)
				return
			}

			err = decoder.Decode(buffer, handler)
			if err != nil {
				log.Error(space, i, " ", err)
				return
			}
		}
	}()

	for i := 0; i < 10; i++ {
		err := encoder.EncodeToWriter(conn, []byte("hello"+powAndString(10, i)))
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func powAndString(x, y int) string {
	r := math.Pow(float64(x), float64(y))
	return strconv.Itoa(int(r))
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
