// github.com/nontechno/later/bar/bar.go

package bar

import "github.com/nontechno/later"

func Example() {
}

func init() {
	later.Register(GetLog, "get.log")
}
