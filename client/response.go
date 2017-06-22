package client

type Response struct {
	Guid  string `json:"guid"`
	Error int32  `json:"error"`
	Value string `json:"value"`
}
