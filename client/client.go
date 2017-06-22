package client

type MerlinClient interface {
	Get(string) (string, error)
	Set(string, string) error
	SetWithExpire(string, string, int64) error
	Delete(string) error
}
