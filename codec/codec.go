package codec

import "io"

//基于长度编码，解决粘包、拆包问题

// 解码器
type Decoder interface {
	Decode(*Buffer, func([]byte)) error
}

// 编码器
type Encoder interface {
	EncodeToWriter(w io.Writer, bytes []byte) error
}
