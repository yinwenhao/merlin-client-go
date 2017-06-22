package handler

type ChannelOutboundHandler interface {
	ChannelHandler
	ChannelWrite(interface{}) (interface{}, error)
}
