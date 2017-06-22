package client

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := NewMerlinClient([]string{"127.0.0.1:5612", "127.0.0.1:5613", "127.0.0.1:5614"})
	if err != nil {
		t.Error(err)
	}
	v, err := client.Get("aaa")
	fmt.Printf("client.Get(\"aaa\") value:%v, error:%v\n", v, err)
	err = client.Set("aaa", "hahaha")
	fmt.Printf("client.Set(\"aaa\", \"hahaha\") error:%v\n", err)
	v, err = client.Get("aaa")
	fmt.Printf("client.Get(\"aaa\") value:%v, error:%v\n", v, err)
	err = client.SetWithExpire("aaa", "kkkkkkkk", 1000)
	fmt.Printf("client.SetWithExpire(\"aaa\", \"kkkkkkkk\") error:%v\n", err)
	v, err = client.Get("aaa")
	fmt.Printf("client.Get(\"aaa\") value:%v, error:%v\n", v, err)
	time.Sleep(time.Second)
	v, err = client.Get("aaa")
	fmt.Printf("client.Get(\"aaa\") after expire value:%v, error:%v\n", v, err)
}
