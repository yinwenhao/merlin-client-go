package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"gitlab.1dmy.com/ezbuy/merlin-client/handler"
)

type SocketTcpClientImpl struct {
	TcpAddr         *net.TCPAddr
	conn            net.Conn
	recvedBuf       *bytes.Buffer
	channelPipeline []handler.ChannelHandler
}

func NewSocketTcpClient(tcpAddr *net.TCPAddr, channelPipeline []handler.ChannelHandler) *SocketTcpClientImpl {
	return &SocketTcpClientImpl{
		TcpAddr:         tcpAddr,
		channelPipeline: channelPipeline,
	}
}

func (m *SocketTcpClientImpl) Close() error {
	return m.conn.Close()
}

func (m *SocketTcpClientImpl) Send(data interface{}) error {
	var o interface{} = data
	for i := len(m.channelPipeline) - 1; i >= 0; i-- {
		if oh, ok := m.channelPipeline[i].(handler.ChannelOutboundHandler); ok {
			oo, err := oh.ChannelWrite(o)
			if err != nil {
				fmt.Printf("doChannelPipeline error:%v", err)
				return err
			}
			o = oo
		}
	}
	v, ok := o.([]byte)
	if !ok {
		return handler.ASSERT_ERROR
	}
	_, err := m.conn.Write(v)
	return err
}

func (m *SocketTcpClientImpl) Connect() {
	conn, err := net.DialTCP("tcp", nil, m.TcpAddr)
	if err != nil {
		panic(err)
	}
	m.conn = conn
	m.recvedBuf = bytes.NewBuffer([]byte{})
	fmt.Println("connect success")
	go m.startListenReceive()
}

func (m *SocketTcpClientImpl) doChannelPipeline(data []byte) {
	var o interface{} = data
	for _, h := range m.channelPipeline {
		if ih, ok := h.(handler.ChannelInboundHandler); ok {
			oo, err := ih.ChannelRead(o)
			if err != nil {
				fmt.Printf("doChannelPipeline error:%v", err)
				return
			}
			o = oo
		}
	}
}

func (m *SocketTcpClientImpl) startListenReceive() {
	var length int32
	var buf [1024]byte
	for {
		n, err := m.conn.Read(buf[0:])
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				fmt.Println("Read timeout")
			} else if err == io.EOF {
				// 连接断开
				fmt.Println("Connection break!")
			} else {
				fmt.Printf("Read from connection failed!, detail:%v\n", err)
			}
			continue
		}
		m.recvedBuf.Write(buf[0:n])
		recvBytes := m.recvedBuf.Len()
		if recvBytes < 4 {
			// 从TCP流中读取数据太少继续读
			fmt.Println("Keep recv...")
			continue
		}

		// 读包头的表示长度的字节
		leaderNumBuf := bytes.NewBuffer(m.recvedBuf.Bytes()[:4])
		binary.Read(leaderNumBuf, binary.BigEndian, &length)
		if recvBytes < int(length)+4 {
			// 从TCP流中读取数据还是太少,继续读
			fmt.Printf("Pack head shows size=%d, buf just recv %d bytes, keep receive.\n", length, recvBytes)
			continue
		}

		if recvBytes > int(length)+4 {
			fmt.Println("Recv data much than a pack.")
		}

		// 读到了完整的包
		go m.doChannelPipeline(m.recvedBuf.Next(int(length) + 4))
	}
}
