package client

type Request struct {
	Guid   string `json:"guid"`
	Expire int64  `json:"expire"`
	Method string `json:"method"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
