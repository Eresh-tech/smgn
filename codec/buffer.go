package codec

import (
	"errors"
	"io"
	"syscall"
)

var ErrorNotEnough = errors.New("Not enough")

// 读缓冲区，每个长连接独占1个
type Buffer struct {
	buf   []byte
	start int //起始位置
	end   int //最后一个字符的下一位，不包含该位置！！！
}

// 新建一个缓冲区
func NewBuffer(bytes []byte) *Buffer {
	return &Buffer{buf: bytes, start: 0, end: 0}
}

// 实际长度
func (b *Buffer) Len() int {
	return b.end - b.start
}

// 容量
func (b *Buffer) Cap() int {
	return len(b.buf)
}

// 取出数据
func (b *Buffer) GetBytes() []byte {
	return b.buf[b.start:b.end]
}

// 获取内部buf的接口
func (b *Buffer) GetBuffer() []byte {
	return b.buf
}

// 重新设置缓存区（将有用字节前移）
func (b *Buffer) reset() {
	if b.start == 0 {
		return
	}

	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
}

// 从文件描述符里读取数据
func (b *Buffer) ReadFromFD(fd int) error {
	b.reset()

	n, err := syscall.Read(fd, b.buf[b.end:])
	if err != nil {
		return err
	}
	if n == 0 {
		return syscall.EAGAIN
	}
	b.end += n
	return nil
}

// 从Reader读取数据
func (b *Buffer) ReadFromReader(reader io.Reader) (int, error) {
	b.reset()
	n, err := reader.Read(b.buf[b.start:b.end])
	if err != nil {
		return n, err
	}
	b.end += n
	return n, nil
}

// 返回n个字节，但不移动文件指针，字节数不足时，返回错误
func (b *Buffer) Seek(len int) ([]byte, error) {
	if b.Len() < len {
		return nil, ErrorNotEnough
	}

	buf := b.buf[b.start : b.start+len]
	return buf, nil
}

// 从start+offset开始读取n个字节，字节数不足时，返回错误
func (b *Buffer) Read(offset, n int) ([]byte, error) {
	if b.Len() < offset+n {
		return nil, ErrorNotEnough
	}
	b.start += offset
	buf := b.buf[b.start : b.start+n]
	b.start += n
	return buf, nil
}

// ReadAll 读取所有字节
func (b *Buffer) ReadAll() []byte {
	buf, _ := b.Read(b.start, b.Len())
	return buf
}
