package handler

import "errors"

type ChannelHandler interface {
}

var ASSERT_ERROR = errors.New("assert failed")
