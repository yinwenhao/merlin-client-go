package client

import (
	"gitlab.1dmy.com/ezbuy/merlin-client/handler"
)

type ResponseHandler struct {
}

var ResponseMap = make(map[string]chan Response)

func (d *ResponseHandler) ChannelRead(data interface{}) (interface{}, error) {
	response, ok := data.(Response)
	if !ok {
		return nil, handler.ASSERT_ERROR
	}
	ResponseMap[response.Guid] <- response
	return response, nil
}
