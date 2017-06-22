package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"gitlab.1dmy.com/ezbuy/merlin-client/handler"
)

type EncodeHandler struct {
}

func (d *EncodeHandler) ChannelWrite(data interface{}) (interface{}, error) {
	r, ok := data.(Request)
	if !ok {
		return nil, handler.ASSERT_ERROR
	}
	message, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return append(append(append([]byte{}, IntToBytes(len(message)+3)...), byte(0), byte(0), byte(0)), message...), nil
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
