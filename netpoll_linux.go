package smgn

import (
	"fmt"
	"golang.org/x/sys/unix"
	"syscall"
)

const (
	EpollRead  = syscall.EPOLLIN | syscall.EPOLLPRI | syscall.EPOLLERR | syscall.EPOLLHUP | unix.EPOLLET | syscall.EPOLLRDHUP
	EpollClose = uint32(syscall.EPOLLIN | syscall.EPOLLRDHUP)
)

type epoll struct {
	listenFD int
	epollFD  int
}

func newNetpoll(address string) (netpoll, error) {
	listenFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	err = syscall.SetsockoptInt(listenFD, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	addr, port, err := getIPPort(address)
	if err != nil {
		return nil, err
	}
	err = syscall.Bind(listenFD, &syscall.SockaddrInet4{
		Port: port,
		Addr: addr,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = syscall.Listen(listenFD, 1024)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &epoll{listenFD: listenFD, epollFD: epollFD}, nil
}

func (n *epoll) accept() (nfd int, addr string, err error) {
	nfd, sa, err := syscall.Accept(n.listenFD)
	if err != nil {
		return
	}

	// 设置为非阻塞状态
	err = syscall.SetNonblock(nfd, true)
	if err != nil {
		return
	}

	err = syscall.EpollCtl(n.epollFD, syscall.EPOLL_CTL_ADD, nfd, &syscall.EpollEvent{
		Events: EpollRead,
		Fd:     int32(nfd),
	})
	if err != nil {
		return
	}

	s := sa.(*syscall.SockaddrInet4)
	addr = fmt.Sprintf("%d.%d.%d.%d:%d", s.Addr[0], s.Addr[1], s.Addr[2], s.Addr[3], s.Port)
	return
}

func (n *epoll) addRead(fd int) error {
	err := syscall.EpollCtl(n.epollFD, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
		Events: EpollRead,
		Fd:     int32(fd),
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *epoll) closeFDRead(fd int) error {
	_, _, e := syscall.Syscall(syscall.SHUT_RD, uintptr(fd), 0, 0)
	if e != 0 {
		return e
	}
	return nil
}
