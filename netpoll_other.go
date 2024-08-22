//go:build !linux
// +build !linux

package smgn

func newNetpoll(address string) (netpoll, error) {
	panic("please run on linux")
}
