package client

import (
	"errors"
	"net"

	"fmt"

	"gitlab.1dmy.com/ezbuy/merlin-client/handler"
	"gitlab.1dmy.com/ezbuy/merlin-client/socket"
	"gitlab.1dmy.com/ezbuy/merlin-client/uuid"
)

const (
	GET    = "get"
	SET    = "set"
	DELETE = "delete"

	NO_EXPIRE = 0
	NO_VALUE  = ""

	DEFAULT_IP_VERSION = "tcp4"
)

type MerlinClientImpl struct {
	scs []socket.SocketClient
	i   int
}

func NewMerlinClient(tcpAddrsString []string) (MerlinClient, error) {
	return NewMerlinClientWithIPVersion(tcpAddrsString, DEFAULT_IP_VERSION)
}

func NewMerlinClientWithIPVersion(tcpAddrsString []string, v string) (MerlinClient, error) {
	if v == "" {
		v = "tcp4"
	}
	tcpAddrs := make([]*net.TCPAddr, len(tcpAddrsString))
	for i, a := range tcpAddrsString {
		tcpAddr, err := net.ResolveTCPAddr(v, a)
		if err != nil {
			return nil, err
		}
		tcpAddrs[i] = tcpAddr
	}
	return NewMerlinClientUseTCPAddrs(tcpAddrs), nil
}

func NewMerlinClientUseTCPAddrs(tcpAddrs []*net.TCPAddr) MerlinClient {
	cp := make([]handler.ChannelHandler, 3)
	cp = append(cp, &DecodeHandler{})
	cp = append(cp, &EncodeHandler{})
	cp = append(cp, &ResponseHandler{})
	scs := make([]socket.SocketClient, len(tcpAddrs))
	for i, tcpAddr := range tcpAddrs {
		scs[i] = socket.NewSocketTcpClient(tcpAddr, cp)
		scs[i].Connect()
	}
	return MerlinClientImpl{
		scs: scs,
		i:   0,
	}
}

func (m MerlinClientImpl) Get(key string) (string, error) {
	r := Request{
		Guid:   getUUID(),
		Expire: NO_EXPIRE,
		Method: GET,
		Key:    key,
		Value:  NO_VALUE,
	}
	res, err := m.sendRequest(r)
	if err != nil {
		return "", err
	}
	if res.Error != 0 {
		return "", fmt.Errorf("gateway error %v", res.Error)
	}
	return res.Value, nil
}

func (m MerlinClientImpl) Set(key string, value string) error {
	return m.SetWithExpire(key, value, NO_EXPIRE)
}

func (m MerlinClientImpl) SetWithExpire(key string, value string, expireMs int64) error {
	r := Request{
		Guid:   getUUID(),
		Expire: expireMs,
		Method: SET,
		Key:    key,
		Value:  value,
	}
	res, err := m.sendRequest(r)
	if err != nil {
		return err
	}
	if res.Error != 0 {
		return fmt.Errorf("gateway error %v", res.Error)
	}
	return err
}

func (m MerlinClientImpl) Delete(key string) error {
	r := Request{
		Guid:   getUUID(),
		Expire: NO_EXPIRE,
		Method: DELETE,
		Key:    key,
		Value:  NO_VALUE,
	}
	res, err := m.sendRequest(r)
	if err != nil {
		return err
	}
	if res.Value != "ok" {
		return errors.New(res.Value)
	}
	return err
}

func (m MerlinClientImpl) sendRequest(request Request) (Response, error) {
	c := make(chan Response, 1)
	ResponseMap[request.Guid] = c
	defer delete(ResponseMap, request.Guid)
	err := m.scs[m.i].Send(request)
	m.i++
	if err != nil {
		return Response{}, err
	}
	return <-ResponseMap[request.Guid], nil
}

func getUUID() string {
	id := uuid.Rand()
	return id.Hex()
}
