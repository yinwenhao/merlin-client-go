package handler

type ChannelInboundHandler interface {
	ChannelHandler
	ChannelRead(interface{}) (interface{}, error)
}
