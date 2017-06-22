package example

import (
	"fmt"

	"gitlab.1dmy.com/ezbuy/merlin-client/client"
)

// 更多例子，请看client包下的test
func example() {
	client, err := client.NewMerlinClient([]string{"127.0.0.1:5612", "127.0.0.1:5613", "127.0.0.1:5614"})
	if err != nil {
		// do something
	}
	v, err := client.Get("aaa")
	fmt.Println(v)
}
