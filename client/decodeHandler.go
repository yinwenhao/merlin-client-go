package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"gitlab.1dmy.com/ezbuy/merlin-client/handler"
)

type DecodeHandler struct {
}

var wrongProtocolVersion = errors.New("wrong protocol version")

func (d *DecodeHandler) ChannelRead(data interface{}) (interface{}, error) {
	message, ok := data.([]byte)
	if !ok {
		return nil, handler.ASSERT_ERROR
	}
	protocolVersion := message[4]
	if protocolVersion != byte(0) {
		return nil, wrongProtocolVersion
	}
	r := &Response{}
	err := json.Unmarshal(message[7:], r)
	if err != nil {
		return nil, err
	}
	return *r, err
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
